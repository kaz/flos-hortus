<template>
  <section id="monitor">
    <h1>Monitors</h1>
    <table>
      <thead>
        <tr>
          <th>Status</th>
          <th>Node</th>
          <th>Monitor</th>
          <th>Timestamp</th>
          <th>Output</th>
        </tr>
      </thead>
      <tbody>
        <tr :key="mon.host+'/'+mon.name" v-for="mon in monitors">
          <td style="text-align: center">{{ mon.status }}</td>
          <td>{{ mon.host }}</td>
          <td>{{ mon.name }}</td>
          <td>{{ mon.timestamp }}</td>
          <td>
            <a v-if="mon.output" href="javascript:" @click="output = mon.output">show</a>
            <span v-else>no output</span>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-show="output">
      <h1>Output</h1>
      <pre>{{ output }}</pre>
      <a href="javascript:" @click="output = null">dissmiss</a>
    </div>
  </section>
</template>

<style>
#monitor table {
  border-collapse: collapse;
}
#monitor th,
#monitor td {
  padding: .5em 2em;
  border: 1px solid #999;
}
#monitor pre {
  padding: 1em;
  background-color: #eee;
}
</style>

<script>
"use strict";

export default {
  data() {
    return {
      timer: null,

      monitors: [],
      output: null,
    }
  },
  created() {
    this.loadMonitors();
    this.timer = setInterval(this.loadMonitors, 10 * 1000);
  },
  beforeDestroy() {
    clearInterval(this.timer);
  },
  methods: {
    async loadMonitors() {
      const resp = await fetch(`/api/database/query`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify("SELECT * FROM lifeline_data"),
      });
      if(resp.status != 200){
        return alert(await resp.text());
      }

      this.monitors = (await resp.json()).rows.map(([host, name, success, output, updated]) => {
        const date = new Date(updated);
        const timestamp = `${date.toLocaleTimeString()}.${date.getMilliseconds()} ${new Date().getTime() - date.getTime() > 30 * 1000 ? " ⚠️" : ""}`;
        const status = parseInt(success) ? "✅" : "❌";
        return {status, host, name, timestamp, output};
      });
    }
  }
}
</script>
