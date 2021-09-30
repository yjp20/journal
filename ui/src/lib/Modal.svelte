<script>
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	function exit(e) {
		if (e.currentTarget != e.target) return;
		dispatch('exit');
	}
</script>

<div class="modal-wrapper" on:click={exit}>
	<div class="modal">
		<slot />
	</div>
</div>

<style>
	.modal-wrapper {
		position: fixed;
		display: flex;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		overflow-y: auto;
		background-color: #00000050;
		z-index: 1000;
	}

	.modal {
		padding: 1em;
		background-color: var(--white);
		border: 1px solid var(--grey);
		border-radius: 5px;
		max-width: 600px;
		flex: 1;
		margin: auto auto;
	}
</style>
