package main

import (
	"context"
	"testing"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// The setup for every test that runs
// similar to changes django testing
func setUp(){
	//use this function for creating a setup necessary for every test
}

// The teardown for the project
func tearDown(){

}

// The function to test the k8 namespace
func TestgetK8Namespace(t * testing.T){
	setUp()

	name, err := getK8Namespace()
	t.Log("The following name was obtained", name)
	if err!=nil{
		t.Error("The K8 namespace was not obtained and the following error occurred:",err)
	}

	tearDown()
}

func TestgetK8RestConfig(t *testing.T){
	setUp()

	config, err := getK8RestConfig()
	t.Log("The following configuration was obtained", config)
	if err!=nil{
		t.Error("The following ")
	}

	tearDown()
}

func TestbecomeLeader(t *testing.T){
	setUp()

	ctx := context.TODO()
	err := becomeLeader(ctx)
	if err!=nil{
		t.Error(" The following error occurred while testing becomeLeader function, the following error occurred",err)
	}

	tearDown()
}

func createManager()(manager.Manager, error){
	name, err := getK8Namespace()
	if err!=nil{
		return nil, err
	}
	config, err := getK8RestConfig()
	if err!=nil{
		return nil,err
	}
	mgr, err := createNewManager(config, name)
	return mgr,err

}
func TestcreateNewManager(t *testing.T){
	setUp()
	mgr, err := createManager()
	if err!=nil{
		t.Error("Creating a manager failed and the following error occurred", err)
	}
	t.Log("The manager obtained",mgr)
	
	tearDown()
}

func TestaddToMetricsPort(t testing.T)  {
	ctx := context.TODO()
	value, err := addToMetricsPort(ctx, metricsPort)
	if err!=nil{
		t.Error("Exposing port failed and the following error occurred", err)
	}
	t.Log("The following value was returned ", value)
}

func TestaddToApiSchema(t testing.T)  {
	mgr, err := createManager()
	if err!=nil{
		t.Error("The following error occurred while creating manager", err)
	}
	t.Log("The value obtained", mgr)
	err = addToApiSchema(mgr)
	if err!=nil{
		t.Error("The following error occurred", err)
	}
	t.Log("The manager was added to api schema")
}

func TestaddToControllerSchema(t testing.T){
	mgr, err := createManager()
	if err!=nil{
		t.Error("The following error occurred", err)
	}
	t.Log("The following value was obtained", mgr)
	err = addToControllerSchema(mgr)
	if err!=nil{
		t.Error("The following error occurred", err)
	}
	t.Log("The manager was added to controller schema")
}