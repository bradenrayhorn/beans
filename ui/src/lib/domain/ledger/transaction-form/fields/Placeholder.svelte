<script lang="ts">
  export let field: string;
  export let value: string;
</script>

<label>
  <span class="label label-text">{field}</span>
  <input
    type="text"
    class="input input-sm input-bordered w-full"
    disabled
    {value}
  />
</label>
