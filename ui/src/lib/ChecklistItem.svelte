<script>
	import Checkbox from '$lib/Checkbox.svelte'
	export let checked
</script>

<div class="checklistitem" class:is-completed={checked} >
	<div class="checklistitem-toggle">
		<Checkbox
			description="Mark as {checked ? 'not done' : 'done'}"
			checked={checked}
			on:change
		/>
	</div>
	<div class="checklistitem-content">
		<div class="checklistitem-description">
			<slot />
		</div>
		<slot name="tags" />
	</div>
	<div class="checklistitem-actions">
		<slot name="actions" />
	</div>
</div>

<style>
	.checklistitem {
		display: flex;
	}

	.checklistitem.is-completed {
		opacity: 0.5;
	}

	.checklistitem-toggle {
		margin-top: 0.1em;
	}

	.checklistitem-content {
		display: block;
		margin-left: 0.25em;
	}

	.checklistitem-actions {
		margin-left: auto;
		align-self: flex-start;
		opacity: 0;
		flex-shrink: 0;
	}

	.checklistitem-description {
		display: inline;
	}

	.checklistitem:hover .checklistitem-description {
		text-decoration: underline;
	}

	.checklistitem:hover .checklistitem-actions {
		opacity: 1;
	}

</style>
