<script setup lang="ts">
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup
import { ref, onMounted } from 'vue'


const WebsocketURL = ref("");

// const WebsocketSchema = ref("ws");
const SchemaDisplayName = ref("ws://");

const ConnectDisplayName = ref("Connect");

let WebSocketConn;

// onMounted(async() => {

// })

const ChangeSchemaAction = async() => {
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

const WebSocketOnMessageHandler = (message) => {
  console.log(message);
}

const WebSocketHandler = async(url: string) => {

  WebSocketConn = NewWebSocket(url);

  WebSocketConn.On
  WebSocketConn.onopen = WebSocketOnOpenHandler;
  WebSocketConn.onclose = WebSocketOnCloseHandler;
  WebSocketConn.onerror = WebSocketOnErrorHandler;
  WebSocketConn.onmessage = WebSocketOnMessageHandler;

  console.log(WebSocketConn);
}

const ConnectAction = async() => {
  ConnectDisplayName.value = "Connecting...";

  let url = SchemaDisplayName.value + WebsocketURL.value;

  await WebSocketHandler(url);
}


</script>

<template>
  <div style="margin: 14px 0">
    <div>
      <el-row class="mb-4">
      <el-input v-model="WebsocketURL" placeholder="Please input websocket url" class="input-with-select">
        <template #prepend>
          <el-button type="success" @click="ChangeSchemaAction">{{ SchemaDisplayName }}</el-button>
<!-- 
          <el-select v-model="WebsocketSchema" placeholder="Select" style="width: 115px">
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
  </div>
</template>

<style scoped>

</style>
