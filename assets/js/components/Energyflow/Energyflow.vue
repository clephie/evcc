<template>
	<div
		class="energyflow cursor-pointer position-relative"
		:class="{ 'energyflow--open': detailsOpen }"
		@click="toggleDetails"
	>
		<div class="row">
			<Visualization
				class="col-12 mb-3 mb-md-4"
				:gridImport="gridImport"
				:selfConsumption="selfConsumption"
				:loadpoints="loadpointsPower"
				:pvExport="pvExport"
				:batteryCharge="batteryCharge"
				:batteryDischarge="batteryDischarge"
				:pvProduction="pvProduction"
				:homePower="homePower"
				:batterySoc="batterySoc"
				:powerInKw="powerInKw"
				:vehicleIcons="vehicleIcons"
			/>
		</div>
		<div class="details" :style="{ height: detailsHeight }">
			<div ref="detailsInner" class="details-inner row">
				<div class="col-12 d-flex justify-content-between pt-2 mb-4">
					<div class="d-flex flex-nowrap align-items-center">
						<span class="color-self me-2"
							><shopicon-filled-square></shopicon-filled-square
						></span>
						<span>{{ $t("main.energyflow.selfConsumption") }}</span>
					</div>
					<div v-if="gridImport > 0" class="d-flex flex-nowrap align-items-center">
						<span>{{ $t("main.energyflow.gridImport") }}</span>
						<span class="color-grid ms-2"
							><shopicon-filled-square></shopicon-filled-square
						></span>
					</div>
					<div v-else class="d-flex flex-nowrap align-items-center">
						<span>{{ $t("main.energyflow.pvExport") }}</span>
						<span class="color-export ms-2"
							><shopicon-filled-square></shopicon-filled-square
						></span>
					</div>
				</div>
				<div
					class="col-12 col-md-6 pe-md-5 pb-4 d-flex flex-column justify-content-between"
				>
					<div class="d-flex justify-content-between align-items-end mb-4">
						<h3 class="m-0">In</h3>
						<span class="fw-bold">
							<AnimatedNumber :to="inPower" :format="kw" />
						</span>
					</div>
					<div>
						<EnergyflowEntry
							:name="$t('main.energyflow.pvProduction')"
							icon="sun"
							:power="pvProduction"
							:powerTooltip="pvTooltip"
							:powerInKw="powerInKw"
						/>
						<EnergyflowEntry
							v-if="batteryConfigured"
							:name="$t('main.energyflow.batteryDischarge')"
							icon="battery"
							:power="batteryDischarge"
							:powerInKw="powerInKw"
							:soc="batterySoc"
							:details="batterySoc"
							:detailsFmt="batteryFmt"
							:detailsTooltip="batteryTooltip"
						/>
						<EnergyflowEntry
							:name="$t('main.energyflow.gridImport')"
							icon="powersupply"
							:power="gridImport"
							:powerInKw="powerInKw"
							:details="detailsValue(tariffGrid, tariffCo2)"
							:detailsFmt="detailsFmt"
							:detailsTooltip="detailsTooltip(tariffGrid, tariffCo2)"
						/>
					</div>
				</div>
				<div
					class="col-12 col-md-6 ps-md-5 pb-4 d-flex flex-column justify-content-between"
				>
					<div class="d-flex justify-content-between align-items-end mb-4">
						<h3 class="m-0">Out</h3>
						<span class="fw-bold">
							<AnimatedNumber :to="outPower" :format="kw" />
						</span>
					</div>
					<div>
						<EnergyflowEntry
							:name="$t('main.energyflow.homePower')"
							icon="home"
							:power="homePower"
							:powerInKw="powerInKw"
							:details="detailsValue(tariffEffectivePrice, tariffEffectiveCo2)"
							:detailsFmt="detailsFmt"
							:detailsTooltip="
								detailsTooltip(tariffEffectivePrice, tariffEffectiveCo2)
							"
						/>
						<EnergyflowEntry
							:name="
								$tc('main.energyflow.loadpoints', activeLoadpointsCount, {
									count: activeLoadpointsCount,
								})
							"
							icon="vehicle"
							:vehicleIcons="vehicleIcons"
							:power="loadpointsPower"
							:powerInKw="powerInKw"
							:details="
								activeLoadpointsCount
									? detailsValue(tariffEffectivePrice, tariffEffectiveCo2)
									: undefined
							"
							:detailsFmt="detailsFmt"
							:detailsTooltip="
								detailsTooltip(tariffEffectivePrice, tariffEffectiveCo2)
							"
						/>
						<EnergyflowEntry
							v-if="batteryConfigured"
							:name="$t('main.energyflow.batteryCharge')"
							icon="battery"
							:power="batteryCharge"
							:powerInKw="powerInKw"
							:soc="batterySoc"
							:details="batterySoc"
							:detailsFmt="batteryFmt"
							:detailsTooltip="batteryTooltip"
						/>
						<EnergyflowEntry
							:name="$t('main.energyflow.pvExport')"
							icon="powersupply"
							:power="pvExport"
							:powerInKw="powerInKw"
							:details="detailsValue(-tariffFeedIn)"
							:detailsFmt="detailsFmt"
							:detailsTooltip="detailsTooltip(-tariffFeedIn)"
						/>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import "@h2d2/shopicons/es/filled/square";
import Visualization from "./Visualization.vue";
import EnergyflowEntry from "./EnergyflowEntry.vue";
import formatter from "../../mixins/formatter";
import AnimatedNumber from "../AnimatedNumber.vue";
import settings from "../../settings";
import { showGridPrice, showGridCo2 } from "../../gridDetails";

export default {
	name: "Energyflow",
	components: { Visualization, EnergyflowEntry, AnimatedNumber },
	mixins: [formatter],
	props: {
		gridConfigured: Boolean,
		gridPower: { type: Number, default: 0 },
		homePower: { type: Number, default: 0 },
		pvConfigured: Boolean,
		pv: { type: Array },
		pvPower: { type: Number, default: 0 },
		loadpointsPower: { type: Number, default: 0 },
		activeLoadpointsCount: { type: Number, default: 0 },
		batteryConfigured: Boolean,
		battery: { type: Array },
		batteryPower: { type: Number, default: 0 },
		batterySoc: { type: Number, default: 0 },
		vehicleIcons: { type: Array },
		tariffGrid: { type: Number },
		tariffFeedIn: { type: Number },
		tariffEffectivePrice: { type: Number },
		tariffCo2: { type: Number },
		tariffEffectiveCo2: { type: Number },
		currency: { type: String },
	},
	data: () => {
		return { detailsOpen: false, detailsCompleteHeight: null };
	},
	computed: {
		gridImport: function () {
			return Math.max(0, this.gridPower);
		},
		pvProduction: function () {
			return Math.abs(this.pvPower);
		},
		batteryPowerAdjusted: function () {
			const batteryPowerThreshold = 50;
			return Math.abs(this.batteryPower) < batteryPowerThreshold ? 0 : this.batteryPower;
		},
		batteryDischarge: function () {
			return Math.abs(Math.max(0, this.batteryPowerAdjusted));
		},
		batteryCharge: function () {
			return Math.abs(Math.min(0, this.batteryPowerAdjusted) * -1);
		},
		selfConsumption: function () {
			const ownPower = this.batteryDischarge + this.pvProduction;
			const consumption = this.homePower + this.batteryCharge + this.loadpointsPower;
			return Math.min(ownPower, consumption);
		},
		pvExport: function () {
			return Math.max(0, this.gridPower * -1);
		},
		powerInKw: function () {
			return this.gridImport + this.selfConsumption + this.pvExport > 1000;
		},
		inPower: function () {
			return this.gridImport + this.pvProduction + this.batteryDischarge;
		},
		outPower: function () {
			return this.homePower + this.loadpointsPower + this.pvExport + this.batteryCharge;
		},
		detailsHeight: function () {
			return this.detailsOpen ? this.detailsCompleteHeight + "px" : 0;
		},
		pvTooltip() {
			if (!Array.isArray(this.pv) || this.pv.length <= 1) {
				return;
			}
			return this.pv.map(({ power }) => this.fmtKw(power, this.powerInKw));
		},
		batteryTooltip() {
			if (!Array.isArray(this.battery)) {
				return;
			}
			return this.battery.map(({ soc, capacity }) => {
				const energy = this.fmtKWh((capacity / 100) * soc * 1e3, true, false, 1);
				const total = this.fmtKWh(capacity * 1e3, true, true, 1);
				return this.$t("main.energyflow.batteryTooltip", { energy, total, soc });
			});
		},
		batteryFmt() {
			return (soc) => `${Math.round(soc)}%`;
		},
	},
	mounted() {
		window.addEventListener("resize", this.updateHeight);
		// height must be calculated in case of initially open details
		if (settings.energyflowDetails) {
			setTimeout(this.toggleDetails, 50);
		}
	},
	unmounted() {
		window.removeEventListener("resize", this.updateHeight);
	},
	methods: {
		detailsTooltip(price, co2) {
			const result = [];
			if (co2 !== undefined) {
				result.push(`${this.fmtCo2Long(co2)}`);
			}
			if (price !== undefined) {
				result.push(`${this.fmtPricePerKWh(price, this.currency)}`);
			}
			return result;
		},
		detailsValue(price, co2) {
			if (showGridPrice() && !isNaN(price)) {
				return price;
			}
			if (showGridCo2() && !isNaN(co2)) {
				return co2;
			}
		},
		detailsFmt(value) {
			if (showGridPrice()) {
				return this.fmtPricePerKWh(value, this.currency, true);
			}
			return this.fmtCo2Short(value);
		},
		kw: function (watt) {
			return this.fmtKw(watt, this.powerInKw);
		},
		toggleDetails: function () {
			this.updateHeight();
			this.detailsOpen = !this.detailsOpen;
			settings.energyflowDetails = this.detailsOpen;
		},
		updateHeight: function () {
			this.detailsCompleteHeight = this.$refs.detailsInner.offsetHeight;
		},
	},
};
</script>
<style scoped>
.details {
	height: 0;
	opacity: 0;
	transform: scale(0.98);
	overflow: visible;
	transition: height, opacity, transform;
	transition-duration: var(--evcc-transition-medium);
	transition-timing-function: cubic-bezier(0.5, 0.5, 0.5, 1.15);
}
.energyflow--open .details {
	opacity: 1;
	transform: scale(1);
}
.color-grid {
	color: var(--evcc-grid);
}
.color-self {
	color: var(--evcc-self);
}
.color-export {
	color: var(--evcc-export);
}
</style>
