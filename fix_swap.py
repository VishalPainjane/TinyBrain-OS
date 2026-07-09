import re

with open('internal/scheduler/v2/backend.cu', 'r') as f:
    content = f.read()

# Fix tb_init
content = re.sub(r'size_t num_layers = 22; // TinyLlama\n    size_t bytes = \(size_t\)config.max_blocks \* num_layers \* 2', r'size_t bytes = (size_t)config.max_blocks * 22 * 2', content)

# Fix tb_swap_out and tb_swap_in
swap_str = '''
tb_error_t tb_swap_out(const char* seq_id, const int* physical_block_ids, int num_blocks) {
    if (g_ctx == nullptr) return TB_ERR_INVALID_STATE;
    std::string sid(seq_id);
    if (g_ctx->swap_registry.find(sid) != g_ctx->swap_registry.end()) {
        return TB_ERR_INVALID_STATE;
    }

    size_t num_layers = 22; // TinyLlama
    size_t block_bytes = num_layers * 2 * g_ctx->block_size * NUM_HEADS * HEAD_DIM * sizeof(half);

    half* host_buffer = nullptr;
    cudaError_t err = cudaMallocHost((void**)&host_buffer, num_blocks * block_bytes);
    if (err != cudaSuccess) {
        return TB_ERR_OOM;
    }
    
    HostSwapAllocation alloc;
    alloc.num_blocks = num_blocks;
    alloc.host_memory = host_buffer;

    for (int i = 0; i < num_blocks; ++i) {
        half* src = g_ctx->kv_cache_vram + ((size_t)physical_block_ids[i] * (block_bytes / sizeof(half)));
        half* dest = host_buffer + ((size_t)i * (block_bytes / sizeof(half)));
        cudaMemcpyAsync(dest, src, block_bytes, cudaMemcpyDeviceToHost, g_ctx->swap_stream);
    }
    
    cudaStreamSynchronize(g_ctx->swap_stream);
    g_ctx->swap_registry[sid] = alloc;
    return TB_SUCCESS;
}

tb_error_t tb_swap_in(const char* seq_id, const int* new_physical_block_ids, int num_blocks) {
    if (g_ctx == nullptr) return TB_ERR_INVALID_STATE;
    std::string sid(seq_id);
    auto it = g_ctx->swap_registry.find(sid);
    if (it == g_ctx->swap_registry.end()) {
        return TB_ERR_INVALID_INPUT;
    }

    HostSwapAllocation alloc = it->second;
    size_t num_layers = 22; // TinyLlama
    size_t block_bytes = num_layers * 2 * g_ctx->block_size * NUM_HEADS * HEAD_DIM * sizeof(half);

    for (int i = 0; i < alloc.num_blocks && i < num_blocks; ++i) {
        half* src = alloc.host_memory + ((size_t)i * (block_bytes / sizeof(half)));
        half* dest = g_ctx->kv_cache_vram + ((size_t)new_physical_block_ids[i] * (block_bytes / sizeof(half)));
        
        cudaMemcpyAsync(dest, src, block_bytes, cudaMemcpyHostToDevice, g_ctx->swap_stream);
    }
    
    cudaStreamSynchronize(g_ctx->swap_stream);
    cudaFreeHost(alloc.host_memory);
    g_ctx->swap_registry.erase(it);
    
    return TB_SUCCESS;
}
'''
content = re.sub(r'tb_error_t tb_swap_out\(const char\* seq_id,.*return TB_SUCCESS;\n}', swap_str.strip(), content, flags=re.DOTALL)

with open('internal/scheduler/v2/backend.cu', 'w') as f:
    f.write(content)
