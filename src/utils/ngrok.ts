import ngrok from "ngrok";
let NGROK_URL = "";

export const InitNgRok = async (port: number): Promise<string> => {
  const url = await ngrok.connect(port);
  NGROK_URL = url;
  return url;
};

export const DisconnectNgRok = async () => {
  await ngrok.disconnect(GetNgRokURL());
};

export const GetNgRokURL = () => {
  return NGROK_URL;
};
