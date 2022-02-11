module github.com/konveyor/tackle-addons-analysis-windup

go 1.16

require (
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/konveyor/tackle-hub v0.0.0-00000000000000-000000000000
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/konveyor/tackle-hub => github.com/mansam/tackle-hub v0.0.0-20220211170149-8180d6dc34c3

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31

replace k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
