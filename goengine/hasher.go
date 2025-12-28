package main

import (
	"goengine/engine"
	"goengine/simple"

	flatbuffers "github.com/google/flatbuffers/go"
)

func handleHashRequest(req *engine.Request) ([]byte, error) {

	inBytes := req.PayloadBytes()
	in := simple.GetRootAsHashRequest(inBytes, 0)
	datalen := in.DataLength()

	inputs := make([]string, datalen)

	for i := 0; i < datalen; i++ {
		inputs[i] = string(in.Data(i))
	}

	b := flatbuffers.NewBuilder(1024)
	offsets := make([]flatbuffers.UOffsetT, len(inputs))

	for i, v := range inputs {
		user := b.CreateString(v)
		hash := b.CreateString("hash_" + v)
		simple.HashResponseStart(b)
		simple.HashResponseAddUserData(b, user)
		simple.HashResponseAddHashData(b, hash)
		offsets[i] = simple.HashResponseEnd(b)
	}

	simple.HashOutPutStartResultVector(b, datalen)
	for i := datalen - 1; i >= 0; i-- {
		b.PrependUOffsetT(offsets[i])
	}

	resultVec := b.EndVector(datalen)

	simple.HashOutPutStart(b)
	simple.HashOutPutAddResult(b, resultVec)
	hashOut := simple.HashOutPutEnd(b)

	b.Finish(hashOut)

	return b.FinishedBytes(), nil
}
