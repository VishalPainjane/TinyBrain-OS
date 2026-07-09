package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
)

type TensorInfo struct {
	Dtype       string `json:"dtype"`
	Shape       []int  `json:"shape"`
	DataOffsets []int  `json:"data_offsets"`
}

func main() {
	numLayers := 2
	hiddenDim := 1024
	numHeads := 8
	headDim := 128
	vocabSize := 32000
	ffnDim := hiddenDim * 4
	_ = numHeads
	_ = headDim

	tensors := make(map[string]TensorInfo)
	tensors["__metadata__"] = TensorInfo{Dtype: "pt"} // dummy metadata to bypass some parsers

	offset := 0

	addTensor := func(name string, shape []int) {
		numElements := 1
		for _, dim := range shape {
			numElements *= dim
		}
		byteSize := numElements * 4
		tensors[name] = TensorInfo{
			Dtype:       "F32",
			Shape:       shape,
			DataOffsets: []int{offset, offset + byteSize},
		}
		offset += byteSize
	}

	addTensor("model.embed_tokens.weight", []int{vocabSize, hiddenDim})

	for i := 0; i < numLayers; i++ {
		prefix := fmt.Sprintf("model.layers.%d.", i)
		addTensor(prefix+"self_attn.q_proj.weight", []int{hiddenDim, hiddenDim})
		addTensor(prefix+"self_attn.k_proj.weight", []int{hiddenDim, hiddenDim})
		addTensor(prefix+"self_attn.v_proj.weight", []int{hiddenDim, hiddenDim})
		addTensor(prefix+"self_attn.o_proj.weight", []int{hiddenDim, hiddenDim})
		addTensor(prefix+"mlp.gate_proj.weight", []int{ffnDim, hiddenDim})
		addTensor(prefix+"mlp.up_proj.weight", []int{ffnDim, hiddenDim})
		addTensor(prefix+"mlp.down_proj.weight", []int{hiddenDim, ffnDim})
	}

	addTensor("lm_head.weight", []int{vocabSize, hiddenDim})

	headerBytes, err := json.Marshal(tensors)
	if err != nil {
		panic(err)
	}

	headerLen := uint64(len(headerBytes))

	outPath := "/tmp/tinybrain_mock.safetensors"
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write 8-byte length
	var lenBuf [8]byte
	binary.LittleEndian.PutUint64(lenBuf[:], headerLen)
	f.Write(lenBuf[:])

	// Write JSON header
	f.Write(headerBytes)

	// Write dummy data (zeros)
	// offset holds total payload size
	fmt.Printf("Total payload size: %d bytes\n", offset)
	chunkSize := 1024 * 1024 * 16 // 16 MB chunk
	zeros := make([]byte, chunkSize)

	written := 0
	for written < offset {
		toWrite := chunkSize
		if offset-written < toWrite {
			toWrite = offset - written
		}
		f.Write(zeros[:toWrite])
		written += toWrite
	}

	fmt.Printf("Successfully generated %s\n", outPath)
}
