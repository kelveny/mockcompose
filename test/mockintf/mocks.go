//go:generate mockcompose -n MockSampleInterface -i SampleInterface
//go:generate mockcompose -n mockFoo -i Foo -p github.com/kelveny/mockcompose/test/foo
//go:generate mockcompose -n mockSecretsInterface -i SecretsInterface -p github.com/kelveny/mockcompose/test/libfn
package mockintf
