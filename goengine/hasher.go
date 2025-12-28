package main

import (
	"crypto/sha256"
	"goengine/engine"
	"goengine/simple"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
)

const MaxWorkers = 64

type job struct {
	index    int
	userData string
}

func handleHashRequest(req *engine.Request) ([]byte, error) {

	inBytes := req.PayloadBytes()
	in := simple.GetRootAsHashRequest(inBytes, 0)
	datalen := in.DataLength()

	inputs := make([]string, datalen)

	for i := range datalen {
		inputs[i] = string(in.Data(i))
	}

	type hashResult struct {
		userdata string
		hash     string
	}

	workerCount := min(len(inputs), MaxWorkers)

	results := make([]hashResult, len(inputs))
	jobs := make(chan job)

	wg := sync.WaitGroup{}

	for range workerCount {
		wg.Go(func() {
			for j := range jobs {
				hashed := buildHash(j.userData)
				results[j.index] = hashResult{
					userdata: j.userData,
					hash:     hashed,
				}
			}
		})
	}

	for i, v := range inputs {
		jobs <- job{
			index:    i,
			userData: v,
		}
	}

	close(jobs)
	wg.Wait()

	b := flatbuffers.NewBuilder(1024)
	offsets := make([]flatbuffers.UOffsetT, len(inputs))

	for i, v := range results {
		user := b.CreateString(v.userdata)
		hash := b.CreateString(v.hash)
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

func buildHash(data string) string {
	// Dummy hash function for illustration
	hashdata := sha256.Sum256([]byte(data))
	return string(hashdata[:])
}
