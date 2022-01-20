import {
  createClient,
  RedisClientType,
  RedisModules,
  RedisScripts,
} from "redis";

let Client: RedisClientType<RedisModules, RedisScripts>;

/**
 * Connect creates a redis client connection
 * @param address Redis Address
 * @returns Promise<void>
 */
export const Connect = async (
  address: string
): Promise<RedisClientType<RedisModules, RedisScripts>> => {
  let client = createClient({
    url: `redis://${address}`,
  });
  Client = client;
  await Client.connect();
  return Client;
};
