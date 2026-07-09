#ifndef TINYBRAIN_BACKEND_H
#define TINYBRAIN_BACKEND_H

#ifdef __cplusplus
extern "C" {
#endif

// tb_error_t defines the strict integer error codes returned by all backend functions,
// ensuring Go can safely catch failure states without dealing with C++ exceptions.
typedef enum {
    TB_SUCCESS = 0,
    TB_ERR_OOM = 1,
    TB_ERR_INVALID_INPUT = 2,
    TB_ERR_CUDA_FAULT = 3
} tb_error_t;

// tb_config_t defines the parameters Go passes during engine startup.
typedef struct {
    int max_batch_size;
    int max_blocks;
    int block_size;     // e.g., 16 tokens per block
    int gpu_id;
    int num_layers;
    int hidden_dim;
    int num_heads;
    int head_dim;
    int vocab_size;
} tb_config_t;

// --- Core Lifecycle Signatures ---

// Boots the resident daemon, allocates the VRAM KV-cache pool, and loads weights.
tb_error_t tb_init(tb_config_t config, const char* model_path);

// Cleans up CUDA streams and frees VRAM.
tb_error_t tb_shutdown();

// --- The Execution Signature ---

// Executes one forward pass iteration for a flattened batch of sequences.
// Note: All input arrays are marked const to guarantee the execution layer
// will not mutate the Go-provided data.
tb_error_t tb_execute_step(
    int num_seqs, 
    const int* token_data, 
    const int* token_offsets, 
    const int* token_lengths, 
    const int* block_data, 
    const int* block_offsets, 
    const int* block_lengths, 
    const int* token_positions,
    float* out_logits
);

// --- The Preemption Signatures (Host RAM Swap) ---

// Instructs CUDA to copy these specific VRAM blocks to Host RAM, keyed by seq_id.
// Go tells C exactly which physical blocks are being evicted.
tb_error_t tb_swap_out(
    const char* seq_id, 
    const int* physical_block_ids, 
    int num_blocks
);

// Instructs CUDA to copy the Host RAM cache for seq_id back into 
// these specific new VRAM block offsets.
tb_error_t tb_swap_in(
    const char* seq_id, 
    const int* new_physical_block_ids, 
    int num_blocks
);

#ifdef __cplusplus
}
#endif

#endif // TINYBRAIN_BACKEND_H
