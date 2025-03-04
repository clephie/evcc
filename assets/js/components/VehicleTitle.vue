<template>
	<div class="d-flex justify-content-between mb-3 align-items-center">
		<h4 class="d-flex align-items-center m-0 flex-grow-1 overflow-hidden">
			<shopicon-regular-refresh
				v-if="icon === 'refresh'"
				ref="refresh"
				data-bs-toggle="tooltip"
				:title="$t('main.vehicle.detectionActive')"
				class="me-2 flex-shrink-0 spin"
			></shopicon-regular-refresh>
			<VehicleIcon
				v-else-if="icon === 'vehicle'"
				:name="vehicleIcon"
				class="me-2 flex-shrink-0"
			/>
			<shopicon-regular-cablecharge
				v-else
				class="me-2 flex-shrink-0"
			></shopicon-regular-cablecharge>
			<VehicleOptions
				v-if="showOptions"
				:id="id"
				class="options"
				:vehicles="otherVehicles"
				:is-unknown="isUnknown"
				@change-vehicle="changeVehicle"
				@remove-vehicle="removeVehicle"
			>
				<span class="flex-grow-1 text-truncate vehicle-name">
					{{ name }}
				</span>
			</VehicleOptions>
			<span v-else class="flex-grow-1 text-truncate vehicle-name">
				{{ name }}
			</span>
		</h4>
	</div>
</template>

<script>
import "@h2d2/shopicons/es/regular/refresh";
import "@h2d2/shopicons/es/regular/cablecharge";
import VehicleIcon from "./VehicleIcon";
import Tooltip from "bootstrap/js/dist/tooltip";

import VehicleOptions from "./VehicleOptions.vue";

export default {
	name: "VehicleTitle",
	components: { VehicleOptions, VehicleIcon },
	props: {
		id: [String, Number],
		vehiclePresent: Boolean,
		vehicleTitle: String,
		vehicleIcon: String,
		vehicleDetectionActive: Boolean,
		parked: Boolean,
		connected: Boolean,
		vehicles: { type: Array, default: () => [] },
	},
	emits: ["change-vehicle", "remove-vehicle"],
	computed: {
		icon() {
			if (this.vehicleDetectionActive) {
				return "refresh";
			}
			if (this.connected || this.parked) {
				return "vehicle";
			}
			return null;
		},
		name() {
			if (this.vehiclePresent || this.parked) {
				return this.vehicleTitle || this.$t("main.vehicle.fallbackName");
			}
			if (this.connected) {
				return this.$t("main.vehicle.unknown");
			}
			return this.$t("main.vehicle.none");
		},
		isUnknown() {
			return !this.vehiclePresent;
		},
		otherVehicles() {
			return this.vehicles
				.map((v, id) => ({
					id: id,
					title: v,
				}))
				.filter((v) => v.title !== this.vehicleTitle);
		},
		showOptions() {
			return !this.isUnknown || this.vehicles.length;
		},
	},
	watch: {
		icon: function () {
			this.tooltip();
		},
	},
	mounted: function () {
		this.tooltip();
	},
	methods: {
		changeVehicle(index) {
			this.$emit("change-vehicle", index);
		},
		removeVehicle() {
			this.$emit("remove-vehicle");
		},
		tooltip() {
			this.$nextTick(() => {
				if (this.$refs.refresh) {
					new Tooltip(this.$refs.refresh);
				}
			});
		},
	},
};
</script>

<style scoped>
.options {
	margin-right: -0.75rem;
}
.vehicle-name {
	text-decoration-color: var(--evcc-gray);
}
.options .vehicle-name {
	text-decoration: underline;
}
.spin {
	animation: rotation 1s infinite cubic-bezier(0.37, 0, 0.63, 1);
}
.spin :deep(svg) {
	/* workaround to fix the not perfectly centered shopicon. Remove once its fixed in @h2d2/shopicons */
	transform: translateY(-0.7px);
}
@keyframes rotation {
	from {
		transform: rotate(0deg);
	}
	to {
		transform: rotate(360deg);
	}
}
</style>
