<script>
</script>

<div class="modal-wrapper">
	<div class="modal">
		<slot />
	</div>
</div>

<style>
	.modal-wrapper {
		position: absolute;
		display: flex;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		overflow-y: auto;
		background-color: #00000050;
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
