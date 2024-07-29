//go:generate mockcompose -n cloneSourceClz -c sourceClz -real Unnamed -real Unnamed2 -real Variadic -real Variadic2 -real Variadic3 -real Variadic4 -real CallFooBar -real CollapsedParams -real CollapsedReturns -real VoidReturn
//go:generate mockcompose -n cloneWithAutoMock -c sourceClz -real "CallPeer,this:.:fmt"
package mockclz
