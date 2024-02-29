<template>
  <q-page class="flex flex-center">
    <q-btn
      color="primary"
      icon="fab fa-google"
      label="Login With Google"
      @click="login"
      v-if="!currentUser"
    />
    <div class="column">
    <p v-if="currentUser" class="text-center">Hello {{ currentUser.email }}</p>
    <q-btn
      color="negative"
      icon="fas fa-sign-out-alt"
      label="Logout"
      @click="logout"
      v-if="currentUser"
    />
    </div>
  </q-page>
</template>

<script>
import axios from "axios";

const apiUrl = "https://idp-example.vercel.app/api";

export default {
  name: "PageIndex",
  data() {
    return {
      currentUser: null
    };
  },
  methods: {
    login() {
      window.location = `${apiUrl}/auth/google/login`;
    },
    logout() {
      window.location = `${apiUrl}/auth/logout`;
    }
  },
  beforeCreate() {
    axios.get(`${apiUrl}/auth/currentUser`).then(resp => {
      if (resp.status == 200) {
        this.currentUser = resp.data;
      }
    }).catch(err => console.log(err));
  }
};
</script>
