<script setup lang="ts">
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup
import { ref, onMounted } from 'vue'


const WebsocketURL = ref("");

// const WebsocketSchema = ref("ws");
const SchemaDisplayName = ref("ws://");

const ConnectDisplayName = ref("Connect");

const SendDisplayName = ref("Send");

const DataType = ref("text");
const InputData = ref("");

const WebsocketTransportData = ref([]);

let WebSocketConn;

// onMounted(async() => {

// })

const ChangeSchemaAction = async () => {
  if (SchemaDisplayName.value.startsWith("ws://")) {
    SchemaDisplayName.value = "wss://";
  } else {
    SchemaDisplayName.value = "ws://";
  }
}

const NewWebSocket = (url: string) => {
  let socket = new WebSocket(url);

  return socket;
}

const WebSocketOnOpenHandler = () => {
  ConnectDisplayName.value = "Connected";
}

const WebSocketOnCloseHandler = () => {
  ConnectDisplayName.value = "Connect";
}

const WebSocketOnErrorHandler = (error) => {
  console.log(error);

  ConnectDisplayName.value = "Connect";
}

const WebSocketOnMessageHandler = (event) => {
  console.log(event);
  console.log('Message from server ', event.data);

  const item = {
    "data": event.data,
    "length": event.data.length,
    "time": event.timeStamp
  };

  WebsocketTransportData.value.push(item);
}

const WebSocketHandler = async (url: string) => {

  WebSocketConn = NewWebSocket(url);

  WebSocketConn.onopen = WebSocketOnOpenHandler;
  WebSocketConn.onclose = WebSocketOnCloseHandler;
  WebSocketConn.onerror = WebSocketOnErrorHandler;
  // WebSocketConn.onmessage = WebSocketOnMessageHandler;

  WebSocketConn.addEventListener('message', WebSocketOnMessageHandler);

  console.log(WebSocketConn);
}

const ConnectAction = async () => {
  ConnectDisplayName.value = "Connecting...";

  let url = SchemaDisplayName.value + WebsocketURL.value;

  await WebSocketHandler(url);
}

const SendAction = async () => {
    if (!WebSocketConn) {
        alert("WebSocketConn is empty");
        return;
    }

    WebSocketConn.send("test");
}


</script>

<template>
  <div class="MainBox">
    <div class="LeftBox">
      <div>
        <el-row class="mb-4">
          <el-input v-model="WebsocketURL" placeholder="Please input websocket url" class="input-with-select">
            <template #prepend>
              <el-button type="success" @click="ChangeSchemaAction">{{ SchemaDisplayName }}</el-button>
              <!--  <el-select v-model="WebsocketSchema" placeholder="Select" style="width: 115px">
            <el-option label="ws://" value="ws" />
            <el-option label="wss://" value="wss" />
          </el-select> -->
            </template>
            <template #append>
              <el-button type="success" @click="ConnectAction">{{ ConnectDisplayName }}</el-button>
            </template>
          </el-input>
        </el-row>
      </div>

      <div style="margin-top: 20px">
        <el-radio-group v-model="DataType">
          <el-radio-button label="text" />
          <el-radio-button label="json" />
          <el-radio-button label="bytes" />
        </el-radio-group>
      </div>
      <div style="margin-top: 20px">
        <el-input class="InputDataBox" v-model="InputData" placeholder="Please input" show-word-limit type="textarea" />
      </div>
      <div style="margin-top: 20px">
        <el-button @click="SendAction">{{ SendDisplayName }}</el-button>
      </div>

    </div>
     <div class="RightBox">
          <el-table :data="WebsocketTransportData" stripe  height="350" style="width: 100%">
            <el-table-column prop="data" label="Data" width="280" />
            <el-table-column prop="length" label="Length" width="80" />
            <el-table-column prop="time" label="Time" />
          </el-table>
     </div>
  </div>
</template>

<style scoped>
.MainBox {
  display:  flex;
  width: 100%;
  margin: 0;
}

.LeftBox {
  margin: 14px 0; flex:  1; border-right: 1px solid #ececec; padding: 20px;
}

.RightBox {
  margin: 14px 0; flex:  2;  padding: 20px;
}

.el-textarea__inner {
  height: 250px !important;
}
</style>
