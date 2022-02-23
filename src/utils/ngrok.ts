import ngrok from "ngrok";
let NGROK_URL = "";

/**
 * InitNgRok starts an ngrok tunnel
 * @param port Port to listen
 * @returns Promise<string> URL of the tunnel
 */
export const InitNgRok = async (port: number): Promise<string> => {
  const url = await ngrok.connect(port);
  NGROK_URL = url;
  return url;
};

/**
 * DisconnectNgRok disconnects ngrok tunnel
 */
export const DisconnectNgRok = async () => {
  await ngrok.disconnect(GetNgRokURL());
};

export const GetNgRokURL = () => {
  return NGROK_URL;
};
