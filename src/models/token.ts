/**
 * Token model
 */
export class Token {
  id: string = "";
  name: string = "";
  decimals: number = 0;
  symbol: string = "";
  address: string = "";

  constructor(token: Token) {
    this.id = token.id || "";
    this.name = token.name || "";
    this.symbol = token.symbol || "";
    this.address = token.address || "";
    this.decimals = token.decimals || 0;
  }
}
