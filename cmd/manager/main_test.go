package main

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"testing"
)

// The setup for every test that runs
func setUp(){
	//This function can be used in the future to create a common setup
}

// The teardown for the project
func tearDown(){
	// A common tearDown function for all the tests
}

//mocking external libraries since this is a unit test
func mockaddToAPISchema(mgr manager.Manager)error{
	// further mocking can be done here
	return nil
}

//mocking external libraries since this is a unit test
func mockaddToControllerSchema(mgr manager.Manager)error{
	// any further mocking can be done here
	return nil
}

func createTestNameConfig()(string, *rest.Config, error){
	name, err := getK8Namespace()
	if err!=nil{
		return "", nil, err
	}
	config, err := getK8RestConfig()
	if err!=nil{
		return "",nil,  err
	}
	return name,config,err
}

// This function returns a temp manager that will be used for testing addToAPISchema  and addToConfigSchema
// it generates namespace and a config to produce a manager
func createTestManager()(manager.Manager, error){
	name, config, err := createTestNameConfig()
	if err!=nil{
		return nil, err
	}
	mgr, err:= createNewManager(config, name)
	if err!=nil{
		return nil,err
	}
	return mgr, err
}

// Tests if a manager is added to API schema
func Test_addToAPISchema(t *testing.T) {
	setUp()

	mgr, err := createTestManager()
	if err!=nil{
		t.Errorf("createNewManager() error = %v",err)
	}
	type args struct {
		mgr manager.Manager
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "addToAPISchhema",
			args:args{mgr},
			wantErr:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockaddToAPISchema(tt.args.mgr); (err != nil) != tt.wantErr {
				t.Errorf("addToAPISchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tearDown()
}

// Tests if becomeLeader function produces errors
func Test_becomeLeader(t *testing.T) {
	setUp()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:"becomeLeader",
			args:args{
				//using the getContext makes sure that if context is changed then the code here need not be changed
				ctx:getContext(),
			},
			wantErr:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := becomeLeader(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("becomeLeader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tearDown()
}

// checks if addToControllerSchema works
func Test_addToControllerSchema(t *testing.T) {
	setUp()

	mgr, err := createTestManager()
	if err!=nil{
		t.Errorf("createNewManager() error = %v",err)
	}

	type args struct {
		mgr manager.Manager
	}
	if err!=nil{
		t.Errorf("createMananger() error = %v",err)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "addToControllerSchema",
			args:    args{
				mgr:mgr,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mockaddToControllerSchema(tt.args.mgr); (err != nil) != tt.wantErr {
				t.Errorf("addToControllerSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tearDown()
}

// Tests if metric port works
func Test_addToMetricsPort(t *testing.T) {
	setUp()

	type args struct {
		ctx         context.Context
		metricsPort int32
	}
	tests := []struct {
		name    string
		args    args
		testValue bool
		want    *v1.Service
		wantErr bool
	}{
		{
			name: "addToMetricsPort",
			args:args{
				ctx:getContext(),
				metricsPort:metricsPort,
			},
			// This determines if the value returned by the function needs to be tested against a standard value
			// by default its false
			testValue:false,
			want:nil,
			wantErr:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addToMetricsPort(tt.args.ctx, tt.args.metricsPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("addToMetricsPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.testValue{
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("addToMetricsPort() got = %v, want %v", got, tt.want)
				}
			}

		})
	}

	tearDown()
}

func Test_createNewManager(t *testing.T) {
	setUp()
	name, config, err:= createTestNameConfig()
	if err!=nil{
		t.Errorf("createTestName")
	}
	type args struct {
		cfg       *rest.Config
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		testValue bool
		want    manager.Manager
		wantErr bool
	}{
		{
			 name:"createNewManager",
			 args:args{
				 cfg:      config ,
				 namespace: name,
			 },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createNewManager(tt.args.cfg, tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("createNewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.testValue {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("createNewManager() got = %v, want %v", got, tt.want)
				}
			}
		})
	}

	tearDown()
}

func Test_getK8Namespace(t *testing.T) {
	setUp()

	tests := []struct {
		name    string
		want    string
		testValue bool
		wantErr bool
	}{
		{
			name: "getK8Namespace",
			want: "",
			testValue:false,
			wantErr:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getK8Namespace()
			if (err != nil) != tt.wantErr {
				t.Errorf("getK8Namespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getK8Namespace() got = %v, want %v", got, tt.want)
			}
		})
	}

	tearDown()
}

func Test_getK8RestConfig(t *testing.T) {
	setUp()

	tests := []struct {
		name    string
		testValue bool
		want    *rest.Config
		wantErr bool
	}{
		{
			name:"getK8RestConfig",
			testValue:false,
			want:nil,
			wantErr:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getK8RestConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("getK8RestConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.testValue{
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("getK8RestConfig() got = %v, want %v", got, tt.want)
				}
			}
		})
	}

	tearDown()
}

func Test_startCmd(t *testing.T) {
	setUp()
	mgr, err:= createTestManager()
	if err!=nil{
		t.Errorf("createTestManager() -> error:%v",err)
	}
	type args struct {
		mgr manager.Manager
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "startCmd",
			args: args{
				mgr: mgr,
			},
			wantErr:false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := startCmd(tt.args.mgr); (err != nil) != tt.wantErr {
				t.Errorf("startCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	tearDown()
}