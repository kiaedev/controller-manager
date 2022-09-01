/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/saltbo/gopkg/strutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	middlewarev1beta1 "my.domain/controller-manager/apis/middleware/v1beta1"
)

// MySQLReconciler reconciles a MySQL object
type MySQLReconciler struct {
	client.Client
	ClientSet *kubernetes.Clientset
	Scheme    *runtime.Scheme
	recorder  record.EventRecorder

	db *sql.DB
}

func NewMySQLReconciler(clientSet *kubernetes.Clientset, client client.Client, scheme *runtime.Scheme, recorder record.EventRecorder, dsn string) *MySQLReconciler {
	db, err := sql.Open("mysql", "root:admin@tcp(localhost:3306)/mysql")
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &MySQLReconciler{
		ClientSet: clientSet,
		Client:    client,
		Scheme:    scheme,
		recorder:  recorder,
		db:        db,
	}
}

//+kubebuilder:rbac:groups=middleware.my.domain,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=middleware.my.domain,resources=mysqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=middleware.my.domain,resources=mysqls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MySQL object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MySQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var rm middlewarev1beta1.MySQL
	if err := r.Get(ctx, req.NamespacedName, &rm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if done, err := r.finalizer(ctx, &rm); err != nil || done {
		return ctrl.Result{}, err
	}

	// generate password or reuse the old password
	password := strutil.RandomText(16)
	pwdSecret := strings.ToLower(fmt.Sprintf("%s-%s", rm.Spec.Database, strutil.RandomText(8)))
	if rm.Status.PwdSecret != "" {
		pwdSecret = rm.Status.PwdSecret
		secret, err := r.ClientSet.CoreV1().Secrets(req.Namespace).Get(ctx, rm.Status.PwdSecret, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			logger.V(3).Info("secret not found, controller will creating new secret by new password")
		} else if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 3}, err
		} else if existPwd, ok := secret.StringData["password"]; ok {
			password = existPwd
		}
	}

	// todo select a mysql instance
	adminDSN := "root:admin@tcp(localhost:3306)/mysql"
	cfg, err := mysql.ParseDSN(adminDSN)
	if err != nil {
		logger.Error(err, "error parsing dsn")
		return ctrl.Result{}, err
	}

	connInfo := middlewarev1beta1.ConnectionInfo{
		Address:   cfg.Addr,
		Database:  rm.Spec.Database,
		Username:  rm.Spec.Username,
		PwdSecret: pwdSecret,
	}
	if err := createOrUpdateDatabase(r.db, connInfo, password); err != nil {
		return ctrl.Result{}, err
	}

	secret := &v1.Secret{}
	secret.SetName(connInfo.PwdSecret)
	secret.SetNamespace(rm.Namespace)
	mutFun := func() (err error) {
		if err := controllerutil.SetOwnerReference(&rm, secret, r.Scheme); err != nil {
			return err
		}

		secret.StringData = map[string]string{
			"password": password,
			"dsn":      connInfo.BuildDSN(password),
		}
		return nil
	}

	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, secret, mutFun); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Get(ctx, req.NamespacedName, &rm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	rm.Status.ConnectionInfo = connInfo
	if err := r.Status().Update(ctx, &rm, &client.UpdateOptions{}); err != nil {
		return ctrl.Result{}, err
	}

	// r.recorder.Event()
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&middlewarev1beta1.MySQL{}).
		Complete(r)
}

func (r *MySQLReconciler) finalizer(ctx context.Context, mysqlDB *middlewarev1beta1.MySQL) (bool, error) {
	myFinalizerName := "mysqls.middleware.my.domain/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if mysqlDB.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(mysqlDB, myFinalizerName) {
			controllerutil.AddFinalizer(mysqlDB, myFinalizerName)
			if err := r.Update(ctx, mysqlDB); err != nil {
				return false, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(mysqlDB, myFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteExternalResources(mysqlDB); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return false, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(mysqlDB, myFinalizerName)
			if err := r.Update(ctx, mysqlDB); err != nil {
				return false, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return true, nil
	}

	return false, nil
}

func (r *MySQLReconciler) deleteExternalResources(mysqlDB *middlewarev1beta1.MySQL) error {
	return deleteDatabase(r.db, mysqlDB.Status.ConnectionInfo)
}

func createOrUpdateDatabase(db *sql.DB, connInfo middlewarev1beta1.ConnectionInfo, password string) error {
	host, _, _ := net.SplitHostPort(connInfo.Address)
	privileges := "Select,Insert,Update,Delete"
	DCLs := []string{
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s ", connInfo.Database),
		fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%s' IDENTIFIED WITH mysql_native_password BY '{%s}'", connInfo.Username, host, password),
		fmt.Sprintf("GRANT %s ON %s.* TO '%s'@'%s'",
			privileges, connInfo.Database, connInfo.Username, host),
	}

	return execDCLs(db, DCLs)
}

func deleteDatabase(db *sql.DB, connInfo middlewarev1beta1.ConnectionInfo) error {
	DCLs := []string{
		fmt.Sprintf("revoke ALL ON `%s`.* FROM '%s'@'%s'", connInfo.Database, connInfo.Username, connInfo.Host()),
		fmt.Sprintf("drop user '%s'@'%s';", connInfo.Username, connInfo.Host()),
		fmt.Sprintf("drop database %s", connInfo.Database),
	}

	return execDCLs(db, DCLs)
}

func reclaimPermission(db *sql.DB, connInfo middlewarev1beta1.ConnectionInfo) error {
	return execDCLs(db, []string{
		fmt.Sprintf("revoke ALL ON `%s`.* FROM '%s'@'%s'", connInfo.Database, connInfo.Username, connInfo.Host()),
	})
}

func execDCLs(db *sql.DB, DCLs []string) error {
	for _, dcl := range DCLs {
		if _, err := db.Exec(dcl); err != nil {
			return nil
		}
	}

	return nil
}
