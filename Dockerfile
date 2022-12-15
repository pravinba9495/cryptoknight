FROM node:16-alpine
RUN mkdir -p /home/pravinba9495/kryptonite
COPY dist /home/pravinba9495/kryptonite
COPY package.json /home/pravinba9495/kryptonite
COPY package-lock.json /home/pravinba9495/kryptonite
WORKDIR /home/pravinba9495/kryptonite
RUN npm install
CMD node index.js