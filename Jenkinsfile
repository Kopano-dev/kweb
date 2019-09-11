#!/usr/bin/env groovy

pipeline {
	agent {
		docker {
			image 'golang:1.13'
			args '-u 0'
		 }
	}
	environment {
		GOBIN = '/usr/local/bin'
		DEBIAN_FRONTEND = 'noninteractive'
		GOLANGCI_LINT_TAG = 'v1.18.0'
	}
	stages {
		stage('Bootstrap') {
			steps {
				echo 'Bootstrapping..'
				sh 'curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(GOBIN)/bin $(GOLANGCI_LINT_TAG)'
				sh 'cd / && go get -v github.com/tebeka/go2xunit'
				sh 'apt-get update && apt-get install -y build-essential'
				sh 'go version'
			}
		}
		stage('Lint') {
			steps {
				echo 'Linting..'
				sh 'make lint-checkstyle'
				checkstyle pattern: 'test/tests.lint.xml', canComputeNew: false, unstableTotalHigh: '100'
			}
		}
		stage('Test') {
			steps {
				echo 'Testing..'
				sh 'make test-xml-short'
				junit allowEmptyResults: true, testResults: 'test/tests.xml'
			}
		}
		stage('Vendor') {
			steps {
				echo 'Fetching vendor dependencies..'
				sh 'make vendor'
			}
		}
		stage('Build') {
			steps {
				echo 'Building..'
				sh 'make DATE=reproducible'
				sh './bin/kwebd version && sha256sum ./bin/kwebd'
			}
		}
		stage('Dist') {
			steps {
				echo 'Dist..'
				sh 'test -z "$(git diff --shortstat 2>/dev/null |tail -n1)" && echo "Clean check passed."'
				sh 'make check'
				sh 'make dist'
			}
		}
	}
	post {
		always {
			archiveArtifacts 'dist/*.tar.gz'
			cleanWs()
		}
	}
}
