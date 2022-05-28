#include "redismodule.h"
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

int RedisModule_OnLoad(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    if (RedisModule_Init(ctx,"avromodule",1,REDISMODULE_APIVER_1)
        == REDISMODULE_ERR) return REDISMODULE_ERR;


    return REDISMODULE_OK;
}

#ifdef __cplusplus
}
#endif
