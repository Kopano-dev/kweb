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
	}
	stages {
		stage('Bootstrap') {
			steps {
				echo 'Bootstrapping..'
				sh 'cd / && go get -v golang.org/x/lint/golint'
				sh 'cd / && go get -v github.com/tebeka/go2xunit'
				sh 'apt-get update && apt-get install -y build-essential'
				sh 'go version'
			}
		}
		stage('Lint') {
			steps {
				echo 'Linting..'
				sh 'make lint 2>&1 | tee golint.txt || true'
				sh 'make vet 2>&1 | tee govet.txt || true'
				warnings parserConfigurations: [[parserName: 'Go Lint', pattern: 'golint.txt'], [parserName: 'Go Vet', pattern: 'govet.txt']], unstableTotalAll: '0'
			}
		}
		stage('Test') {
			steps {
				echo 'Testing..'
				sh 'make test-xml-short'
				junit allowEmptyResults: true, testResults: 'test/*.xml'
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
