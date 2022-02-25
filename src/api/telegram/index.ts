import Axios from "axios";

export class Telegram {
  static SendMessage = async (
    token: string,
    chatId: string,
    message: string,
    notify: boolean = true
  ): Promise<any> => {
    return new Promise((resolve, reject) => {
      Axios.post(
        `https://api.telegram.org/bot${token}/sendMessage`,
        {
          chat_id: chatId,
          text: message,
          disable_notification: !notify,
        },
        {
          timeout: 5000,
        }
      )
        .then((response) => {
          resolve(response);
        })
        .catch((error) => {
          reject(error);
        });
    });
  };

  static SetWebhook = async (token: string, url: string): Promise<any> => {
    return new Promise((resolve, reject) => {
      Axios.post(
        `https://api.telegram.org/bot${token}/setWebhook`,
        {
          url,
          allowed_updates: ["message"],
          drop_pending_updates: true,
        },
        {
          timeout: 5000,
        }
      )
        .then((response) => {
          if (response.data && response.data.ok) {
            resolve(response.data);
          } else {
            reject(response.data);
          }
        })
        .catch((error) => {
          reject(error);
        });
    });
  };
}

export default Telegram;
