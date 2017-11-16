<template>
  <div class="d-flex align-items-end">
    <div class="flex-1-auto mw-100">
      <p>{{contact.username}}</p>
    </div>
  </div>
</template>

<script>
  export default {
     props: ['contact']
  }
</script>
