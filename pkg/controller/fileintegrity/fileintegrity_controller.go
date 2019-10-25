package fileintegrity

import (
	"context"
	"errors"

	appsv1 "k8s.io/api/apps/v1"

	fileintegrityv1alpha1 "github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1"
	"github.com/mrogers950/file-integrity-operator/pkg/common"

	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_fileintegrity")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new FileIntegrity Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFileIntegrity{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("fileintegrity-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource FileIntegrity
	err = c.Watch(&source.Kind{Type: &fileintegrityv1alpha1.FileIntegrity{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner FileIntegrity
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &fileintegrityv1alpha1.FileIntegrity{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileFileIntegrity implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFileIntegrity{}

// ReconcileFileIntegrity reconciles a FileIntegrity object
type ReconcileFileIntegrity struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a FileIntegrity object and makes changes based on the state read
// and what is in the FileIntegrity.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileFileIntegrity) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("reconciling FileIntegrity")

	// Fetch the FileIntegrity instance
	instance := &fileintegrityv1alpha1.FileIntegrity{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if kerr.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// create the aide script configmap if it does not exist.
	defaultAideScript := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: common.AideScriptConfigMapName, Namespace: common.FileIntegrityNamespace}, defaultAideScript)
	if err != nil {
		if !kerr.IsNotFound(err) {
			reqLogger.Error(err, "error getting aide script")
			return reconcile.Result{}, err
		}
		// does not exist, create it
		createErr := r.client.Create(context.TODO(), defaultAIDEScript())
		if createErr != nil {
			reqLogger.Error(err, "error creating aide script")
			return reconcile.Result{}, createErr
		}
	}

	// handle configuration configmap
	defaultAideConf := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: common.DefaultConfigMapName, Namespace: common.FileIntegrityNamespace}, defaultAideConf)
	if err != nil {
		if !kerr.IsNotFound(err) {
			reqLogger.Error(err, "error getting default aide config")
			return reconcile.Result{}, err
		}
		// does not exist, create it
		createErr := r.client.Create(context.TODO(), defaultAIDEConfigMap())
		if createErr != nil {
			reqLogger.Error(err, "error creating default aide config")
			return reconcile.Result{}, createErr
		}
	}
	if _, ok := defaultAideConf.Data[common.DefaultConfDataKey]; !ok {
		reqLogger.Info("default aide.conf has no data")
		return reconcile.Result{}, errors.New("default aide.conf has no data")
	}

	// handle user-provided configmap
	defaultAideConfCopy := defaultAideConf.DeepCopy()
	reqLogger.Info("instance spec", "Instance.Spec", instance.Spec)
	if len(instance.Spec.Config.Name) > 0 && len(instance.Spec.Config.Namespace) > 0 {
		reqLogger.Info("checking for configmap update")

		cm := &corev1.ConfigMap{}
		cfErr := r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.Config.Name, Namespace: instance.Spec.Config.Namespace}, cm)
		if cfErr != nil {
			if !kerr.IsNotFound(cfErr) {
				reqLogger.Error(cfErr, "error getting aide config configmap")
				return reconcile.Result{}, cfErr
			}
		}
		reqLogger.Info("default configmap found")
		if !kerr.IsNotFound(cfErr) {
			key := common.DefaultConfDataKey
			if instance.Spec.Config.Key != "" {
				key = instance.Spec.Config.Key
			}
			conf, ok := cm.Data[key]
			if ok && len(conf) > 0 && conf != defaultAideConfCopy.Data[common.DefaultConfDataKey] {
				reqLogger.Info("preparing aide conf")
				preparedConf, prepErr := prepareAideConf(conf)
				if prepErr != nil {
					reqLogger.Error(prepErr, "error preparing provided aide conf")
					return reconcile.Result{}, prepErr
				}
				reqLogger.Info("updating aide conf")
				defaultAideConfCopy.Data[common.DefaultConfDataKey] = preparedConf
				updateErr := r.client.Update(context.TODO(), defaultAideConfCopy)
				if updateErr != nil {
					reqLogger.Error(updateErr, "error updating default configmap")
					return reconcile.Result{}, updateErr
				}
			}
		}
	}

	reqLogger.Info("reconciling daemonSets")
	daemonSet := &appsv1.DaemonSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: common.WorkerDaemonSetName, Namespace: common.FileIntegrityNamespace}, daemonSet)
	if err != nil {
		if !kerr.IsNotFound(err) {
			reqLogger.Error(err, "error getting worker daemonSet")
			return reconcile.Result{}, err
		}
		// create
		ds := workerAideDaemonset()
		createErr := r.client.Create(context.TODO(), ds)
		if createErr != nil {
			reqLogger.Error(createErr, "error creating worker daemonSet")
			return reconcile.Result{}, createErr
		}
	}
	masterDaemonSet := &appsv1.DaemonSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: common.MasterDaemonSetName, Namespace: common.FileIntegrityNamespace}, masterDaemonSet)
	if err != nil {
		if !kerr.IsNotFound(err) {
			reqLogger.Error(err, "error getting master daemonSet")
			return reconcile.Result{}, err
		}
		mds := masterAideDaemonset()
		mcreateErr := r.client.Create(context.TODO(), mds)
		if mcreateErr != nil {
			reqLogger.Error(mcreateErr, "error creating master daemonSet")
			return reconcile.Result{}, mcreateErr
		}
	}
	return reconcile.Result{}, nil
}

func defaultAIDEConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.DefaultConfigMapName,
			Namespace: common.FileIntegrityNamespace,
		},
		Data: map[string]string{
			"aide.conf": defaultAideConfig,
		},
	}
}

func defaultAIDEScript() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.AideScriptConfigMapName,
			Namespace: common.FileIntegrityNamespace,
		},
		Data: map[string]string{
			"aide.sh": aideScript,
		},
	}
}

func workerAideDaemonset() *appsv1.DaemonSet {
	priv := true
	runAs := int64(0)
	mode := int32(0744)

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.WorkerDaemonSetName,
			Namespace: common.FileIntegrityNamespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": common.WorkerDaemonSetName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": common.WorkerDaemonSetName,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: common.OperatorServiceAccountName,
					Containers: []corev1.Container{
						{
							SecurityContext: &corev1.SecurityContext{
								Privileged: &priv,
								RunAsUser:  &runAs,
							},
							Name:    "aide",
							Image:   "docker.io/mrogers950/aide:latest",
							Command: []string{common.AideScriptPath},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "hostroot",
									MountPath: "/hostroot",
								},
								{
									Name:      "config",
									MountPath: "/tmp",
								},
								{
									Name:      common.AideScriptConfigMapName,
									MountPath: "/scripts",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "hostroot",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/",
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: common.DefaultConfigMapName,
									},
								},
							},
						},
						{
							Name: common.AideScriptConfigMapName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: common.AideScriptConfigMapName,
									},
									DefaultMode: &mode,
								},
							},
						},
					},
				},
			},
		},
	}
}

func masterAideDaemonset() *appsv1.DaemonSet {
	priv := true
	runAs := int64(0)
	mode := int32(0744)

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      common.MasterDaemonSetName,
			Namespace: common.FileIntegrityNamespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": common.MasterDaemonSetName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": common.MasterDaemonSetName,
					},
				},
				Spec: corev1.PodSpec{
					Tolerations: []corev1.Toleration{
						{
							Key:      "node-role.kubernetes.io/master",
							Operator: "Exists",
							Effect:   "NoSchedule",
						},
					},
					NodeSelector: map[string]string{
						"node-role.kubernetes.io/master": "",
					},
					ServiceAccountName: common.OperatorServiceAccountName,
					Containers: []corev1.Container{
						{
							SecurityContext: &corev1.SecurityContext{
								Privileged: &priv,
								RunAsUser:  &runAs,
							},
							Name:    "aide",
							Image:   "docker.io/mrogers950/aide:latest",
							Command: []string{common.AideScriptPath},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "hostroot",
									MountPath: "/hostroot",
								},
								{
									Name:      "config",
									MountPath: "/tmp",
								},
								{
									Name:      common.AideScriptConfigMapName,
									MountPath: "/scripts",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "hostroot",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/",
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "aide-conf",
									},
								},
							},
						},
						{
							Name: common.AideScriptConfigMapName,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: common.AideScriptConfigMapName,
									},
									DefaultMode: &mode,
								},
							},
						},
					},
				},
			},
		},
	}
}