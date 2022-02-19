import Axios from "axios";

/**
 * SendMessage sends a message to a client on telegram
 * @param token Telegram Bot Token
 * @param chatId CHat ID of the receiver
 * @param message Message to send
 * @returns Promise<any>
 */
export const SendMessage = async (
  token: string,
  chatId: string,
  message: string
): Promise<any> => {
  return new Promise((resolve, reject) => {
    Axios.post(`https://api.telegram.org/bot${token}/sendMessage`, {
      chat_id: chatId,
      text: message,
    })
      .then((response) => {
        resolve(response);
      })
      .catch((error) => {
        reject(error);
      });
  });
};
