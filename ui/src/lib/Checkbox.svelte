<script>
	export let checked;
	export let description = undefined;
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
	<div class="checkbox-checkmark">✓</div>
</label>

<style>
	.checkbox {
		position: relative;
		width: 0.9em;
		height: 0.9em;
		color: inherit;
		border: 1px solid currentColor;
		border-radius: 999px;
		box-sizing: border-box;
	}

	.checkbox:hover {
		background-color: #f9f9f9;
	}

	.checkbox-input {
		display: none;
	}

	.checkbox-checkmark {
		opacity: 0;
		font-size: 0.6em;
		text-align: center;
	}

	:checked + .checkbox-checkmark {
		opacity: 1;
	}

	.checkbox:hover .checkbox-checkmark {
		opacity: 0.5;
	}
</style>
