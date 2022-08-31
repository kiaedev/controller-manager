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

	"github.com/go-sql-driver/mysql"
	"github.com/saltbo/gopkg/strutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	middlewarev1beta1 "my.domain/controller-manager/apis/middleware/v1beta1"
)

// MySQLReconciler reconciles a MySQL object
type MySQLReconciler struct {
	client.Client
	ClientSet *kubernetes.Clientset
	Scheme    *runtime.Scheme
	recorder  record.EventRecorder

	dsn string
}

func NewMySQLReconciler(clientSet *kubernetes.Clientset, client client.Client, scheme *runtime.Scheme, recorder record.EventRecorder, dsn string) *MySQLReconciler {
	return &MySQLReconciler{
		ClientSet: clientSet,
		Client:    client,
		Scheme:    scheme,
		recorder:  recorder,
		dsn:       dsn,
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
	_ = log.FromContext(ctx)

	var res_mysql middlewarev1beta1.MySQL
	if err := r.Get(ctx, req.NamespacedName, &res_mysql); err != nil {
		return ctrl.Result{}, err
	}

	password := strutil.RandomText(16)
	if err := createDatabase(res_mysql.Name, password); err != nil {
		return ctrl.Result{}, err
	}

	secret := v1.Secret("", res_mysql.Namespace)
	secret.WithStringData(map[string]string{
		"password": password,
	})

	_, err := r.ClientSet.CoreV1().Secrets(res_mysql.Namespace).Apply(ctx, secret, metav1.ApplyOptions{})
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Update(ctx, &res_mysql, &client.UpdateOptions{}); err != nil {
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

func createDatabase(name, password string) error {
	dsn := ""
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	cfg, _ := mysql.ParseDSN(dsn)
	database := name
	username := fmt.Sprintf("%s_rw", name)
	privileges := "Select,Insert,Update,Delete"

	sqlFormats := []string{
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database),
		fmt.Sprintf("CREATE USER IF NOT EXISTS %s", username),
		fmt.Sprintf("GRANT %s ON %s.* TO '%s'@'%s' identified by '%s'", privileges, database, username, cfg.Addr, password),
	}

	for _, format := range sqlFormats {
		_, err = db.Exec(fmt.Sprintf(format))
		if err != nil {
			return err
		}
	}

	return nil
}
