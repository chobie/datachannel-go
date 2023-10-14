# datachannel-go 

This is a cgo binding for libdatachannel.

# Purpose

Currently, the main goals are to playing the CAPI implementation of libdatachannel and provide sample code. 
Not intended for use in commercial environments.

# How to build

* Ensure that the dylib of libdatachannel is already installed.
* Run `cd examples/simple_datachannel_client && go build`.

If you encounter issues with dynamic linking resolution at runtime, add the module installation path using:

`export DYLD_LIBRARY_PATH=/usr/local/lib:$DYLD_LIBRARY_PATH`.

# Current Status

I have probably implemented all of CAPI, but I haven't tested its functionality.
At the least, the code in the example works.

# License

MPL 2.0