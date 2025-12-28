from cffi import FFI
import engine.Request as Req
import flatbuffers
from simple import HashRequest as In
from simple import HashOutPut as Out

ffi = FFI()
ffi.cdef("""
int CallEngine(void* inPtr, unsigned long long inLen,
               unsigned char** outPtr, unsigned long long* outLen);
void FreeBuffer(void* ptr);
""")


import os

lib_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "bin", "goengine.dll"))
C = ffi.dlopen(lib_path)

data = ["dip", "raj"]
b = flatbuffers.Builder(1024)

# --------------------------
# building HashRequest
# --------------------------

hashData_offsets = [b.CreateString(u.decode('utf-8') if isinstance(u, bytes) else str(u)) for u in data]
In.HashRequestStartDataVector(b, len(data))
for u in reversed(hashData_offsets):
        b.PrependUOffsetTRelative(u)
hashData_vec = b.EndVector(len(hashData_offsets))

In.HashRequestStart(b)
In.HashRequestAddData(b, hashData_vec)
In.HashRequestAddCount(b, len(data))

payload = In.HashRequestEnd(b)
b.Finish(payload)

payload_bytes = b.Output()

# --------------------------
# Wrapping into engine Request
# --------------------------


b2 = flatbuffers.Builder(1024)
payload_vec = b2.CreateByteVector(payload_bytes)

Req.RequestStart(b2)
Req.RequestAddFuncId(b2, 1)
Req.RequestAddRequestId(b2, 1)
Req.RequestAddPayload(b2, payload_vec)
req = Req.RequestEnd(b2)
b2.Finish(req)

req_bytes = b2.Output()


out_ptr = ffi.new("unsigned char**")
out_len = ffi.new("unsigned long long*")

rc = C.CallEngine(
    ffi.from_buffer(req_bytes),
    len(req_bytes),
    out_ptr,
    out_len,
)

print("CallEngine rc=", rc)


if rc != 0:
    raise RuntimeError("CallEngine failed")

out_buf = ffi.buffer(out_ptr[0], out_len[0])[:]
C.FreeBuffer(out_ptr[0])

out = Out.HashOutPut.GetRootAsHashOutPut(out_buf, 0)

for i in range(out.ResultLength()):
    o= out.Result(i)
    print(o.HashData().decode('utf-8'))
