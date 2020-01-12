# Development Workflow

## Chaos-Operator Architecture

![Chaos-Operator Architecture](./docs/chaos-operator-architecture.png)

## Prerequisites

* You have Go 1.10+ installed on your local host/development machine.
* You have Docker installed on your local host/development machine. Docker is required for building chaos-operator component container images and to push them into a Kubernetes cluster for testing.

## Initial Setup

### Fork in the cloud

1. Visit https://github.com/litmuschaos/chaos-operator.
2. Click `Fork` button (top right) to establish a cloud-based fork.

### Clone fork to local host

Place `litmuschaos/chaos-operator` code on your `GOPATH` using the following cloning procedure.
Create your clone:

```sh

mkdir -p $GOPATH/src/github.com/litmuschaos
cd $GOPATH/src/github.com/litmuschaos

# Note: Here $user is your GitHub profile name
git clone https://github.com/$user/chaos-operator.git

# Configure remote upstream
cd $GOPATH/src/github.com/litmuschaos/chaos-operator
git remote add upstream https://github.com/litmuschaos/chaos-operator.git

# Never push to upstream master
git remote set-url --push upstream no_push

# Confirm that your remotes make sense
git remote -v
```

## Development

### Always sync your local repository

Open a terminal on your local host. Change directory to the chaos-operator fork root.

```sh
$ cd $GOPATH/src/github.com/litmuschaos/chaos-operator
```

 Checkout the master branch.

 ```sh
 $ git checkout master
 Switched to branch 'master'
 Your branch is up-to-date with 'origin/master'.
 ```

 Recall that origin/master is a branch on your remote GitHub repository.
 Make sure you have the upstream remote litmuschaos/chaos-operator by listing them.

 ```sh
 $ git remote -v
 origin	https://github.com/$user/chaos-operator.git (fetch)
 origin	https://github.com/$user/chaos-operator.git (push)
 upstream	https://github.com/litmuschaos/chaos-operator.git (fetch)
 upstream	no_push (push)
 ```

 If the upstream is missing, add it by using below command.

 ```sh
 $ git remote add upstream https://github.com/litmuschaos/chaos-operator.git
 ```

 Fetch all the changes from the upstream master branch.

 ```sh
 $ git fetch upstream master
 remote: Counting objects: 141, done.
 remote: Compressing objects: 100% (29/29), done.
 remote: Total 141 (delta 52), reused 46 (delta 46), pack-reused 66
 Receiving objects: 100% (141/141), 112.43 KiB | 0 bytes/s, done.
 Resolving deltas: 100% (79/79), done.
 From github.com:litmuschaos/chaos-operator
   * branch            master     -> FETCH_HEAD
 ```

 Rebase your local master with the upstream/master.

 ```sh
 $ git rebase upstream/master
 First, rewinding head to replay your work on top of it...
 Fast-forwarded master to upstream/master.
 ```

 This command applies all the commits from the upstream master to your local master.

 Check the status of your local branch.

 ```sh
 $ git status
 On branch master
 Your branch is ahead of 'origin/master' by 38 commits.
 (use "git push" to publish your local commits)
 nothing to commit, working directory clean
 ```

 Your local repository now has all the changes from the upstream remote. You need to push the changes to your own remote fork which is origin master.

 Push the rebased master to origin master.

 ```sh
 $ git push origin master
 Username for 'https://github.com': $user
 Password for 'https://$user@github.com':
 Counting objects: 223, done.
 Compressing objects: 100% (38/38), done.
 Writing objects: 100% (69/69), 8.76 KiB | 0 bytes/s, done.
 Total 69 (delta 53), reused 47 (delta 31)
 To https://github.com/$user/chaos-operator.git
 8e107a9..5035fa1  master -> master
 ```

### Create a new feature branch to work on your issue

 Your branch name should have the format `XYZ-descriptive` where `XYZ` is the issue number you are working on followed by some descriptive text. For example:

 ```sh
 $ git checkout -b 256-fix-reconsiler
 Switched to a new branch '256-reconsiler'
 ```

### Make your changes and build them

 ```sh
 cd $GOPATH/src/github.com/litmuschaos/chaos-operator
 make all
 ```

Check your linting.

 ```sh
 make lint
 ```

### Test your changes

- Replace the image with the builded image [here](../deploy/operator.yaml)

- Run the choos-operator in kubernetes cluster
 ```sh
 cd $GOPATH/src/github.com/litmuschaos/chaos-operator
 kubectl apply -f ./deploy/chaos_crds.yaml
 kubectl apply -f ./deploy/rbac.yaml
 kubectl apply -f ./deploy/operator.yaml
 ```
- Run the chaos by following the [Litmus Docs](https://docs.litmuschaos.io/docs/getstarted/#install-chaos-experiments)

- Verify the changes

### Keep your branch in sync

[Rebasing](https://git-scm.com/docs/git-rebase) is very important to keep your branch in sync with the changes being made by others and to avoid huge merge conflicts while raising your Pull Requests. You will always have to rebase before raising the PR.

```sh
git fetch upstream
git rebase upstream/master
```

While you rebase your changes, you must resolve any conflicts that might arise and build and test your changes using the above steps.

## Submission

### Create a pull request

Before you raise the Pull Requests, ensure you have reviewed the checklist in the [CONTRIBUTING GUIDE](../CONTRIBUTING.md):
- Ensure that you have re-based your changes with the upstream using the steps above.
- Ensure that you have added the required unit tests for the bug fixes or new feature that you have introduced.
- Ensure your commits history is clean with proper header and descriptions.

Go to the [litmuschaos/chaos-operator github](https://github.com/litmuschaos/chaos-operator) and follow the Open Pull Request link to raise your PR from your development branch.
