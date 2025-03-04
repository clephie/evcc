<template>
	<Teleport to="body">
		<div
			id="globalSettingsModal"
			class="modal fade text-dark"
			data-bs-backdrop="true"
			tabindex="-1"
			role="dialog"
			aria-hidden="true"
		>
			<div class="modal-dialog modal-dialog-centered" role="document">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">{{ $t("settings.title") }}</h5>
						<button
							type="button"
							class="btn-close"
							data-bs-dismiss="modal"
							aria-label="Close"
						></button>
					</div>
					<div class="modal-body">
						<div class="container">
							<FormRow id="settingsDesign" :label="$t('settings.theme.label')">
								<SelectGroup
									id="settingsDesign"
									v-model="theme"
									class="w-100"
									:options="
										THEMES.map((value) => ({
											value,
											name: $t(`settings.theme.${value}`),
										}))
									"
								/>
							</FormRow>
							<FormRow id="settingsLanguage" :label="$t('settings.language.label')">
								<select
									id="settingsLanguage"
									v-model="language"
									class="form-select form-select-sm w-75"
								>
									<option value="">{{ $t("settings.language.auto") }}</option>
									<option
										v-for="option in languageOptions"
										:key="option"
										:value="option.value"
									>
										{{ option.name }}
									</option>
								</select>
							</FormRow>
							<FormRow id="settingsUnit" :label="$t('settings.unit.label')">
								<SelectGroup
									id="settingsUnit"
									v-model="unit"
									class="w-75"
									:options="
										UNITS.map((value) => ({
											value,
											name: $t(`settings.unit.${value}`),
										}))
									"
								/>
							</FormRow>
							<FormRow
								id="settingsGridDetails"
								:label="$t('settings.gridDetails.label')"
							>
								<SelectGroup
									id="settingsGridDetails"
									v-model="gridDetails"
									class="w-100"
									:options="
										GRID_DETAILS.map((value) => ({
											value,
											name: $t(`settings.gridDetails.${value}`),
											disabled: isDisabled(value),
										}))
									"
								/>
							</FormRow>
							<FormRow id="telemetryEnabled" :label="$t('settings.telemetry.label')">
								<TelemetrySettings :sponsor="sponsor" class="mt-1 mb-0" />
							</FormRow>
							<FormRow
								v-if="isNightly"
								id="hiddenFeaturesEnabled"
								:label="`${$t('settings.hiddenFeatures.label')} 🧪`"
							>
								<div class="form-check form-switch my-1">
									<input
										id="hiddenFeaturesEnabled"
										v-model="hiddenFeatures"
										class="form-check-input"
										type="checkbox"
										role="switch"
									/>
									<div class="form-check-label">
										<label for="telemetryEnabled">
											{{ $t("settings.hiddenFeatures.value") }}
										</label>
									</div>
								</div>
							</FormRow>
						</div>
					</div>
				</div>
			</div>
		</div>
	</Teleport>
</template>

<script>
import TelemetrySettings from "./TelemetrySettings.vue";
import FormRow from "./FormRow.vue";
import SelectGroup from "./SelectGroup.vue";
import { getLocalePreference, setLocalePreference, LOCALES, removeLocalePreference } from "../i18n";
import { getThemePreference, setThemePreference, THEMES } from "../theme";
import { getUnits, setUnits, UNITS } from "../units";
import { getGridDetails, setGridDetails, GRID_DETAILS } from "../gridDetails";
import { getHiddenFeatures, setHiddenFeatures } from "../featureflags";

export default {
	name: "GlobalSettingsModal",
	components: { TelemetrySettings, FormRow, SelectGroup },
	props: {
		sponsor: String,
		hasPrice: Boolean,
		hasCo2: Boolean,
	},
	data: function () {
		return {
			theme: getThemePreference(),
			language: getLocalePreference() || "",
			unit: getUnits(),
			gridDetails: getGridDetails(),
			hiddenFeatures: getHiddenFeatures(),
			THEMES,
			UNITS,
			GRID_DETAILS,
		};
	},
	computed: {
		languageOptions: () => {
			const locales = Object.entries(LOCALES).map(([key, value]) => {
				return { value: key, name: value[1] };
			});
			// sort by name
			locales.sort((a, b) => (a.name < b.name ? -1 : 1));
			return locales;
		},
		isNightly: () => {
			return !!window.evcc.commit;
		},
	},
	watch: {
		unit(value) {
			setUnits(value);
		},
		gridDetails(value) {
			setGridDetails(value);
		},
		theme(value) {
			setThemePreference(value);
		},
		hiddenFeatures(value) {
			setHiddenFeatures(value);
		},
		language(value) {
			const i18n = this.$root.$i18n;
			if (value) {
				setLocalePreference(i18n, value);
			} else {
				removeLocalePreference(i18n);
			}
		},
	},
	methods: {
		isDisabled(option) {
			return (option === "co2" && !this.hasCo2) || (option === "price" && !this.hasPrice);
		},
	},
};
</script>
<style scoped>
.container {
	margin-left: calc(var(--bs-gutter-x) * -0.5);
	margin-right: calc(var(--bs-gutter-x) * -0.5);
}
</style>
