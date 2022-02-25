import {
  createClient,
  RedisClientType,
  RedisModules,
  RedisScripts,
} from "redis";

export const NewClient = async (
  address: string
): Promise<RedisClientType<RedisModules, RedisScripts>> => {
  let client = createClient({
    url: `redis://${address}`,
  });
  await client.connect();
  return client;
};
