<script>
	export let checked
	export let description = undefined
</script>

<label class="checkbox">
	<input
		aria-label={description}
		class="checkbox-input"
		type="checkbox"
		on:change
		on:click
		bind:checked
	/>
	<div class="checkbox-mark"></div>
</label>

<style>
	.checkbox-mark {
		position: relative;
		display: inline-block;
		width: 0.9rem;
		height: 0.9rem;
		color: inherit;
		outline: 2px solid white;
		outline-offset: -3px;
		border: 1px solid currentColor;
		border-radius: 4px;
		box-sizing: border-box;
		transition: 0.1s border-width ease;
	}

	.checkbox:hover > .checkbox-mark {
		border-width: 0.45rem;
	}

	:checked + .checkbox-mark {
		border-width: 0.45rem;
	}

	.checkbox-input {
		display: none;
	}
</style>
