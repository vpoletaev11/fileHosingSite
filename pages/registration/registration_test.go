package registration

import (
	"database/sql/driver"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type anyString struct{}

// ()Match() checks is cookie value valid
func (a anyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	if !ok {
		return false
	}
	if !(len(v.(string)) == 60) {
		return false
	}
	return true
}

// TestPageSuccessGET checks workability of GET requests handler in Page()
func TestPageSuccessGET(t *testing.T) {
	sut := Page(nil)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/registration", nil)
	require.NoError(t, err)

	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageSuccessPost checks workability of POST requests handler in Page()
func TestPageSuccessPost(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("INSERT INTO users").WithArgs("example", anyString{}, "Europe/Moscow").WillReturnResult(sqlmock.NewResult(1, 1))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")
	data.Add("timezone", "Europe/Moscow")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	assert.Equal(t, "/login", w.Header().Get("Location"))
}

// TestPageMissingTemplate tests case when template file is missing.
// Cannot be runned in parallel.
func TestPageMissingTemplate(t *testing.T) {
	// renaming exists template file
	oldName := "../../" + pathTemplateRegistration
	newName := "../../" + pathTemplateRegistration + "edit"
	err := os.Rename(oldName, newName)
	require.NoError(t, err)
	lenOrigName := len(oldName)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://localhost/registration", nil)
	require.NoError(t, err)

	// running of the page handler with un-exists template file
	sut := Page(nil)
	sut(w, r)

	assert.Equal(t, 500, w.Code)

	// renaming template file to original filename
	defer func() {
		// renaming template file to original filename
		oldName = newName
		newName = oldName[:lenOrigName]
		err = os.Rename(oldName, newName)
		require.NoError(t, err)
	}()

	// checking error handler works correct
	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
}

// TestPageEmptyUsername tests case when username is empty.
func TestPageEmptyUsername(t *testing.T) {
	data := url.Values{}
	data.Set("username", "")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Username cannot be empty</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageEmptyPassword1 tests case when password1 is empty.
func TestPageEmptyPassword1(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Password cannot be empty</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageEmptyPassword2 tests case when password2 is empty.
func TestPageEmptyPassword2(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Password cannot be empty</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageLargerUsername tests case when len(username) > 20.
func TestPageLargerUsername(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example_larger_than_20_characters")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Username cannot be longer than 20 characters</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageLargerPassword1 tests case when len(password1) > 40.
func TestPageLargerPassword1(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example_larger_than_40_characters_example_larger_than_40_characters")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Password cannot be longer than 40 characters</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageLargerPassword2 tests case when len(password2) > 40.
func TestPageLargerPassword2(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example_larger_than_40_characters_example_larger_than_40_characters")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Password cannot be longer than 40 characters</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageNonLowerCaseUsername tests case when username is non lower-case
func TestPageNonLowerCaseUsernam(t *testing.T) {
	data := url.Values{}
	data.Set("username", "Example")
	data.Add("password1", "example")
	data.Add("password2", "example")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Please use lower case username</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageMismatchingPasswords tests case when password1 != password2
func TestPageMismatchingPasswords(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example1")
	data.Add("password2", "example2")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Passwords doesn't match</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageNotUniqueUsername tests case when username already exists in MySQL database
func TestPageNotUniqueUsername(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("INSERT INTO users").WithArgs("example", anyString{}, "Europe/Moscow").WillReturnError(fmt.Errorf("Error 1062"))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")
	data.Add("timezone", "Europe/Moscow")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Username already used</h2>
            </form>
        </div>
    </body>`, bodyString)
}

// TestPageUsernameInsertionDBInternalError tests case when username insertion in MySQL database unreachable because of internal error.
func TestPageUsernameInsertionDBInternalError(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	sqlMock.ExpectExec("INSERT INTO users").WithArgs("example", anyString{}, "Europe/Moscow").WillReturnError(fmt.Errorf("Testing error"))

	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")
	data.Add("timezone", "Europe/Moscow")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(db)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">INTERNAL ERROR. Please try later</h2>
            </form>
        </div>
    </body>`, bodyString)
}

func TestPageNonSelectedTimezoneError(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")
	data.Add("timezone", "empty")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Please set your timezone</h2>
            </form>
        </div>
    </body>`, bodyString)
}

func TestPageEmptyTimezoneError(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")
	data.Add("timezone", "")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Incorrect timezone</h2>
            </form>
        </div>
    </body>`, bodyString)
}

func TestPageWrongTimezoneError(t *testing.T) {
	data := url.Values{}
	data.Set("username", "example")
	data.Add("password1", "example")
	data.Add("password2", "example")
	data.Add("timezone", "wrong timezone")

	r, err := http.NewRequest("POST", "http://localhost/registration", strings.NewReader(data.Encode()))
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	w := httptest.NewRecorder()

	sut := Page(nil)
	sut(w, r)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, `<!doctype html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Registration</title>
    <link rel="stylesheet" href="assets/css/register.css">

    <head>

    <body bgcolor=#f1ded3>
        <div class="registerForm">
            <form action="" method="post">
                <p>Create username: <input required type="text" name="username"></p>
                <p>Create password: <input required type="password" name="password1"></p>
                <p>Repeat password: <input required type="password" name="password2"></p>
                <p>Set Timezone: <select name="timezone">
                    <option value="empty">Choose timezone</option>
                    <optgroup label="Africa">
                        <option value="Africa/Abidjan">(+00:00 UTC) Abidjan</option>
                        <option value="Africa/Accra">(+00:00 UTC) Accra</option>
                        <option value="Africa/Bamako">(+00:00 UTC) Bamako</option>
                        <option value="Africa/Banjul">(+00:00 UTC) Banjul</option>
                        <option value="Africa/Bissau">(+00:00 UTC) Bissau</option>
                        <option value="Africa/Conakry">(+00:00 UTC) Conakry</option>
                        <option value="Africa/Dakar">(+00:00 UTC) Dakar</option>
                        <option value="Africa/Freetown">(+00:00 UTC) Freetown</option>
                        <option value="Africa/Lome">(+00:00 UTC) Lome</option>
                        <option value="Africa/Monrovia">(+00:00 UTC) Monrovia</option>
                        <option value="Africa/Nouakchott">(+00:00 UTC) Nouakchott</option>
                        <option value="Africa/Ouagadougou">(+00:00 UTC) Ouagadougou</option>
                        <option value="Africa/Algiers">(+01:00 UTC) Algiers</option>
                        <option value="Africa/Bangui">(+01:00 UTC) Bangui</option>
                        <option value="Africa/Brazzaville">(+01:00 UTC) Brazzaville</option>
                        <option value="Africa/Casablanca">(+01:00 UTC) Casablanca</option>
                        <option value="Africa/Ceuta">(+01:00 UTC) Ceuta</option>
                        <option value="Africa/Douala">(+01:00 UTC) Douala</option>
                        <option value="Africa/El_Aaiun">(+01:00 UTC) El Aaiun</option>
                        <option value="Africa/Kinshasa">(+01:00 UTC) Kinshasa</option>
                        <option value="Africa/Lagos">(+01:00 UTC) Lagos</option>
                        <option value="Africa/Libreville">(+01:00 UTC) Libreville</option>
                        <option value="Africa/Luanda">(+01:00 UTC) Luanda</option>
                        <option value="Africa/Malabo">(+01:00 UTC) Malabo</option>
                        <option value="Africa/Ndjamena">(+01:00 UTC) Ndjamena</option>
                        <option value="Africa/Niamey">(+01:00 UTC) Niamey</option>
                        <option value="Africa/Porto-Novo">(+01:00 UTC) Porto-Novo</option>
                        <option value="Africa/Sao_Tome">(+01:00 UTC) Sao Tome</option>
                        <option value="Africa/Tunis">(+01:00 UTC) Tunis</option>
                        <option value="Africa/Blantyre">(+02:00 UTC) Blantyre</option>
                        <option value="Africa/Bujumbura">(+02:00 UTC) Bujumbura</option>
                        <option value="Africa/Cairo">(+02:00 UTC) Cairo</option>
                        <option value="Africa/Gaborone">(+02:00 UTC) Gaborone</option>
                        <option value="Africa/Harare">(+02:00 UTC) Harare</option>
                        <option value="Africa/Johannesburg">(+02:00 UTC) Johannesburg</option>
                        <option value="Africa/Khartoum">(+02:00 UTC) Khartoum</option>
                        <option value="Africa/Kigali">(+02:00 UTC) Kigali</option>
                        <option value="Africa/Lubumbashi">(+02:00 UTC) Lubumbashi</option>
                        <option value="Africa/Lusaka">(+02:00 UTC) Lusaka</option>
                        <option value="Africa/Maputo">(+02:00 UTC) Maputo</option>
                        <option value="Africa/Maseru">(+02:00 UTC) Maseru</option>
                        <option value="Africa/Mbabane">(+02:00 UTC) Mbabane</option>
                        <option value="Africa/Tripoli">(+02:00 UTC) Tripoli</option>
                        <option value="Africa/Windhoek">(+02:00 UTC) Windhoek</option>
                        <option value="Africa/Addis_Ababa">(+03:00 UTC) Addis Ababa</option>
                        <option value="Africa/Asmara">(+03:00 UTC) Asmara</option>
                        <option value="Africa/Dar_es_Salaam">(+03:00 UTC) Dar es Salaam</option>
                        <option value="Africa/Djibouti">(+03:00 UTC) Djibouti</option>
                        <option value="Africa/Juba">(+03:00 UTC) Juba</option>
                        <option value="Africa/Kampala">(+03:00 UTC) Kampala</option>
                        <option value="Africa/Mogadishu">(+03:00 UTC) Mogadishu</option>
                        <option value="Africa/Nairobi">(+03:00 UTC) Nairobi</option>
                    </optgroup>
                    <optgroup label="America">
                        <option value="America/Adak">(-10:00 UTC) Adak</option>
                        <option value="America/Anchorage">(-09:00 UTC) Anchorage</option>
                        <option value="America/Juneau">(-09:00 UTC) Juneau</option>
                        <option value="America/Metlakatla">(-09:00 UTC) Metlakatla</option>
                        <option value="America/Nome">(-09:00 UTC) Nome</option>
                        <option value="America/Sitka">(-09:00 UTC) Sitka</option>
                        <option value="America/Yakutat">(-09:00 UTC) Yakutat</option>
                        <option value="America/Dawson">(-08:00 UTC) Dawson</option>
                        <option value="America/Los_Angeles">(-08:00 UTC) Los Angeles</option>
                        <option value="America/Tijuana">(-08:00 UTC) Tijuana</option>
                        <option value="America/Vancouver">(-08:00 UTC) Vancouver</option>
                        <option value="America/Whitehorse">(-08:00 UTC) Whitehorse</option>
                        <option value="America/Boise">(-07:00 UTC) Boise</option>
                        <option value="America/Cambridge_Bay">(-07:00 UTC) Cambridge Bay</option>
                        <option value="America/Chihuahua">(-07:00 UTC) Chihuahua</option>
                        <option value="America/Creston">(-07:00 UTC) Creston</option>
                        <option value="America/Dawson_Creek">(-07:00 UTC) Dawson Creek</option>
                        <option value="America/Denver">(-07:00 UTC) Denver</option>
                        <option value="America/Edmonton">(-07:00 UTC) Edmonton</option>
                        <option value="America/Fort_Nelson">(-07:00 UTC) Fort Nelson</option>
                        <option value="America/Hermosillo">(-07:00 UTC) Hermosillo</option>
                        <option value="America/Inuvik">(-07:00 UTC) Inuvik</option>
                        <option value="America/Mazatlan">(-07:00 UTC) Mazatlan</option>
                        <option value="America/Ojinaga">(-07:00 UTC) Ojinaga</option>
                        <option value="America/Phoenix">(-07:00 UTC) Phoenix</option>
                        <option value="America/Yellowknife">(-07:00 UTC) Yellowknife</option>
                        <option value="America/Bahia_Banderas">(-06:00 UTC) Bahia Banderas</option>
                        <option value="America/Belize">(-06:00 UTC) Belize</option>
                        <option value="America/Chicago">(-06:00 UTC) Chicago</option>
                        <option value="America/Costa_Rica">(-06:00 UTC) Costa Rica</option>
                        <option value="America/El_Salvador">(-06:00 UTC) El Salvador</option>
                        <option value="America/Guatemala">(-06:00 UTC) Guatemala</option>
                        <option value="America/Indiana/Knox">(-06:00 UTC) Indiana/Knox</option>
                        <option value="America/Indiana/Tell_City">(-06:00 UTC) Indiana/Tell City</option>
                        <option value="America/Managua">(-06:00 UTC) Managua</option>
                        <option value="America/Matamoros">(-06:00 UTC) Matamoros</option>
                        <option value="America/Menominee">(-06:00 UTC) Menominee</option>
                        <option value="America/Merida">(-06:00 UTC) Merida</option>
                        <option value="America/Mexico_City">(-06:00 UTC) Mexico City</option>
                        <option value="America/Monterrey">(-06:00 UTC) Monterrey</option>
                        <option value="America/North_Dakota/Beulah">(-06:00 UTC) North Dakota/Beulah</option>
                        <option value="America/North_Dakota/Center">(-06:00 UTC) North Dakota/Center</option>
                        <option value="America/North_Dakota/New_Salem">(-06:00 UTC) North Dakota/New Salem</option>
                        <option value="America/Rainy_River">(-06:00 UTC) Rainy River</option>
                        <option value="America/Rankin_Inlet">(-06:00 UTC) Rankin Inlet</option>
                        <option value="America/Regina">(-06:00 UTC) Regina</option>
                        <option value="America/Resolute">(-06:00 UTC) Resolute</option>
                        <option value="America/Swift_Current">(-06:00 UTC) Swift Current</option>
                        <option value="America/Tegucigalpa">(-06:00 UTC) Tegucigalpa</option>
                        <option value="America/Winnipeg">(-06:00 UTC) Winnipeg</option>
                        <option value="America/Atikokan">(-05:00 UTC) Atikokan</option>
                        <option value="America/Bogota">(-05:00 UTC) Bogota</option>
                        <option value="America/Cancun">(-05:00 UTC) Cancun</option>
                        <option value="America/Cayman">(-05:00 UTC) Cayman</option>
                        <option value="America/Detroit">(-05:00 UTC) Detroit</option>
                        <option value="America/Eirunepe">(-05:00 UTC) Eirunepe</option>
                        <option value="America/Grand_Turk">(-05:00 UTC) Grand Turk</option>
                        <option value="America/Guayaquil">(-05:00 UTC) Guayaquil</option>
                        <option value="America/Havana">(-05:00 UTC) Havana</option>
                        <option value="America/Indiana/Indianapolis">(-05:00 UTC) Indiana/Indianapolis</option>
                        <option value="America/Indiana/Marengo">(-05:00 UTC) Indiana/Marengo</option>
                        <option value="America/Indiana/Petersburg">(-05:00 UTC) Indiana/Petersburg</option>
                        <option value="America/Indiana/Vevay">(-05:00 UTC) Indiana/Vevay</option>
                        <option value="America/Indiana/Vincennes">(-05:00 UTC) Indiana/Vincennes</option>
                        <option value="America/Indiana/Winamac">(-05:00 UTC) Indiana/Winamac</option>
                        <option value="America/Iqaluit">(-05:00 UTC) Iqaluit</option>
                        <option value="America/Jamaica">(-05:00 UTC) Jamaica</option>
                        <option value="America/Kentucky/Louisville">(-05:00 UTC) Kentucky/Louisville</option>
                        <option value="America/Kentucky/Monticello">(-05:00 UTC) Kentucky/Monticello</option>
                        <option value="America/Lima">(-05:00 UTC) Lima</option>
                        <option value="America/Nassau">(-05:00 UTC) Nassau</option>
                        <option value="America/New_York">(-05:00 UTC) New York</option>
                        <option value="America/Nipigon">(-05:00 UTC) Nipigon</option>
                        <option value="America/Panama">(-05:00 UTC) Panama</option>
                        <option value="America/Pangnirtung">(-05:00 UTC) Pangnirtung</option>
                        <option value="America/Port-au-Prince">(-05:00 UTC) Port-au-Prince</option>
                        <option value="America/Rio_Branco">(-05:00 UTC) Rio Branco</option>
                        <option value="America/Thunder_Bay">(-05:00 UTC) Thunder Bay</option>
                        <option value="America/Toronto">(-05:00 UTC) Toronto</option>
                        <option value="America/Anguilla">(-04:00 UTC) Anguilla</option>
                        <option value="America/Antigua">(-04:00 UTC) Antigua</option>
                        <option value="America/Aruba">(-04:00 UTC) Aruba</option>
                        <option value="America/Barbados">(-04:00 UTC) Barbados</option>
                        <option value="America/Blanc-Sablon">(-04:00 UTC) Blanc-Sablon</option>
                        <option value="America/Boa_Vista">(-04:00 UTC) Boa Vista</option>
                        <option value="America/Caracas">(-04:00 UTC) Caracas</option>
                        <option value="America/Curacao">(-04:00 UTC) Curacao</option>
                        <option value="America/Dominica">(-04:00 UTC) Dominica</option>
                        <option value="America/Glace_Bay">(-04:00 UTC) Glace Bay</option>
                        <option value="America/Goose_Bay">(-04:00 UTC) Goose Bay</option>
                        <option value="America/Grenada">(-04:00 UTC) Grenada</option>
                        <option value="America/Guadeloupe">(-04:00 UTC) Guadeloupe</option>
                        <option value="America/Guyana">(-04:00 UTC) Guyana</option>
                        <option value="America/Halifax">(-04:00 UTC) Halifax</option>
                        <option value="America/Kralendijk">(-04:00 UTC) Kralendijk</option>
                        <option value="America/La_Paz">(-04:00 UTC) La Paz</option>
                        <option value="America/Lower_Princes">(-04:00 UTC) Lower Princes</option>
                        <option value="America/Manaus">(-04:00 UTC) Manaus</option>
                        <option value="America/Marigot">(-04:00 UTC) Marigot</option>
                        <option value="America/Martinique">(-04:00 UTC) Martinique</option>
                        <option value="America/Moncton">(-04:00 UTC) Moncton</option>
                        <option value="America/Montserrat">(-04:00 UTC) Montserrat</option>
                        <option value="America/Port_of_Spain">(-04:00 UTC) Port of Spain</option>
                        <option value="America/Porto_Velho">(-04:00 UTC) Porto Velho</option>
                        <option value="America/Puerto_Rico">(-04:00 UTC) Puerto Rico</option>
                        <option value="America/Santo_Domingo">(-04:00 UTC) Santo Domingo</option>
                        <option value="America/St_Barthelemy">(-04:00 UTC) St Barthelemy</option>
                        <option value="America/St_Kitts">(-04:00 UTC) St Kitts</option>
                        <option value="America/St_Lucia">(-04:00 UTC) St Lucia</option>
                        <option value="America/St_Thomas">(-04:00 UTC) St Thomas</option>
                        <option value="America/St_Vincent">(-04:00 UTC) St Vincent</option>
                        <option value="America/Thule">(-04:00 UTC) Thule</option>
                        <option value="America/Tortola">(-04:00 UTC) Tortola</option>
                        <option value="America/Araguaina">(-03:00 UTC) Araguaina</option>
                        <option value="America/Argentina/Buenos_Aires">(-03:00 UTC) Argentina/Buenos Aires</option>
                        <option value="America/Argentina/Catamarca">(-03:00 UTC) Argentina/Catamarca</option>
                        <option value="America/Argentina/Cordoba">(-03:00 UTC) Argentina/Cordoba</option>
                        <option value="America/Argentina/Jujuy">(-03:00 UTC) Argentina/Jujuy</option>
                        <option value="America/Argentina/La_Rioja">(-03:00 UTC) Argentina/La Rioja</option>
                        <option value="America/Argentina/Mendoza">(-03:00 UTC) Argentina/Mendoza</option>
                        <option value="America/Argentina/Rio_Gallegos">(-03:00 UTC) Argentina/Rio Gallegos</option>
                        <option value="America/Argentina/Salta">(-03:00 UTC) Argentina/Salta</option>
                        <option value="America/Argentina/San_Juan">(-03:00 UTC) Argentina/San Juan</option>
                        <option value="America/Argentina/San_Luis">(-03:00 UTC) Argentina/San Luis</option>
                        <option value="America/Argentina/Tucuman">(-03:00 UTC) Argentina/Tucuman</option>
                        <option value="America/Argentina/Ushuaia">(-03:00 UTC) Argentina/Ushuaia</option>
                        <option value="America/Asuncion">(-03:00 UTC) Asuncion</option>
                        <option value="America/Bahia">(-03:00 UTC) Bahia</option>
                        <option value="America/Belem">(-03:00 UTC) Belem</option>
                        <option value="America/Campo_Grande">(-03:00 UTC) Campo Grande</option>
                        <option value="America/Cayenne">(-03:00 UTC) Cayenne</option>
                        <option value="America/Cuiaba">(-03:00 UTC) Cuiaba</option>
                        <option value="America/Fortaleza">(-03:00 UTC) Fortaleza</option>
                        <option value="America/Godthab">(-03:00 UTC) Godthab</option>
                        <option value="America/Maceio">(-03:00 UTC) Maceio</option>
                        <option value="America/Miquelon">(-03:00 UTC) Miquelon</option>
                        <option value="America/Montevideo">(-03:00 UTC) Montevideo</option>
                        <option value="America/Paramaribo">(-03:00 UTC) Paramaribo</option>
                        <option value="America/Punta_Arenas">(-03:00 UTC) Punta Arenas</option>
                        <option value="America/Recife">(-03:00 UTC) Recife</option>
                        <option value="America/Santarem">(-03:00 UTC) Santarem</option>
                        <option value="America/Santiago">(-03:00 UTC) Santiago</option>
                        <option value="America/St_Johns">(-03:30 UTC) St Johns</option>
                        <option value="America/Noronha">(-02:00 UTC) Noronha</option>
                        <option value="America/Sao_Paulo">(-02:00 UTC) Sao Paulo</option>
                        <option value="America/Scoresbysund">(-01:00 UTC) Scoresbysund</option>
                        <option value="America/Danmarkshavn">(+00:00 UTC) Danmarkshavn</option>
                    </optgroup>
                    <optgroup label="Antarctica">
                        <option value="Antarctica/Palmer">(-03:00 UTC) Palmer</option>
                        <option value="Antarctica/Rothera">(-03:00 UTC) Rothera</option>
                        <option value="Antarctica/Troll">(+00:00 UTC) Troll</option>
                        <option value="Antarctica/Syowa">(+03:00 UTC) Syowa</option>
                        <option value="Antarctica/Mawson">(+05:00 UTC) Mawson</option>
                        <option value="Antarctica/Vostok">(+06:00 UTC) Vostok</option>
                        <option value="Antarctica/Davis">(+07:00 UTC) Davis</option>
                        <option value="Antarctica/Casey">(+08:00 UTC) Casey</option>
                        <option value="Antarctica/DumontDUrville">(+10:00 UTC) DumontDUrville</option>
                        <option value="Antarctica/Macquarie">(+11:00 UTC) Macquarie</option>
                        <option value="Antarctica/McMurdo">(+13:00 UTC) McMurdo</option>
                    </optgroup>
                    <optgroup label="Asia">
                        <option value="Asia/Amman">(+02:00 UTC) Amman</option>
                        <option value="Asia/Beirut">(+02:00 UTC) Beirut</option>
                        <option value="Asia/Damascus">(+02:00 UTC) Damascus</option>
                        <option value="Asia/Famagusta">(+02:00 UTC) Famagusta</option>
                        <option value="Asia/Gaza">(+02:00 UTC) Gaza</option>
                        <option value="Asia/Hebron">(+02:00 UTC) Hebron</option>
                        <option value="Asia/Jerusalem">(+02:00 UTC) Jerusalem</option>
                        <option value="Asia/Nicosia">(+02:00 UTC) Nicosia</option>
                        <option value="Asia/Aden">(+03:00 UTC) Aden</option>
                        <option value="Asia/Baghdad">(+03:00 UTC) Baghdad</option>
                        <option value="Asia/Bahrain">(+03:00 UTC) Bahrain</option>
                        <option value="Asia/Kuwait">(+03:00 UTC) Kuwait</option>
                        <option value="Asia/Qatar">(+03:00 UTC) Qatar</option>
                        <option value="Asia/Riyadh">(+03:00 UTC) Riyadh</option>
                        <option value="Asia/Tehran">(+03:30 UTC) Tehran</option>
                        <option value="Asia/Baku">(+04:00 UTC) Baku</option>
                        <option value="Asia/Dubai">(+04:00 UTC) Dubai</option>
                        <option value="Asia/Muscat">(+04:00 UTC) Muscat</option>
                        <option value="Asia/Tbilisi">(+04:00 UTC) Tbilisi</option>
                        <option value="Asia/Yerevan">(+04:00 UTC) Yerevan</option>
                        <option value="Asia/Kabul">(+04:30 UTC) Kabul</option>
                        <option value="Asia/Aqtau">(+05:00 UTC) Aqtau</option>
                        <option value="Asia/Aqtobe">(+05:00 UTC) Aqtobe</option>
                        <option value="Asia/Ashgabat">(+05:00 UTC) Ashgabat</option>
                        <option value="Asia/Atyrau">(+05:00 UTC) Atyrau</option>
                        <option value="Asia/Dushanbe">(+05:00 UTC) Dushanbe</option>
                        <option value="Asia/Karachi">(+05:00 UTC) Karachi</option>
                        <option value="Asia/Oral">(+05:00 UTC) Oral</option>
                        <option value="Asia/Samarkand">(+05:00 UTC) Samarkand</option>
                        <option value="Asia/Tashkent">(+05:00 UTC) Tashkent</option>
                        <option value="Asia/Yekaterinburg">(+05:00 UTC) Yekaterinburg</option>
                        <option value="Asia/Colombo">(+05:30 UTC) Colombo</option>
                        <option value="Asia/Kolkata">(+05:30 UTC) Kolkata</option>
                        <option value="Asia/Kathmandu">(+05:45 UTC) Kathmandu</option>
                        <option value="Asia/Almaty">(+06:00 UTC) Almaty</option>
                        <option value="Asia/Bishkek">(+06:00 UTC) Bishkek</option>
                        <option value="Asia/Dhaka">(+06:00 UTC) Dhaka</option>
                        <option value="Asia/Omsk">(+06:00 UTC) Omsk</option>
                        <option value="Asia/Qyzylorda">(+06:00 UTC) Qyzylorda</option>
                        <option value="Asia/Thimphu">(+06:00 UTC) Thimphu</option>
                        <option value="Asia/Urumqi">(+06:00 UTC) Urumqi</option>
                        <option value="Asia/Yangon">(+06:30 UTC) Yangon</option>
                        <option value="Asia/Bangkok">(+07:00 UTC) Bangkok</option>
                        <option value="Asia/Barnaul">(+07:00 UTC) Barnaul</option>
                        <option value="Asia/Ho_Chi_Minh">(+07:00 UTC) Ho Chi Minh</option>
                        <option value="Asia/Hovd">(+07:00 UTC) Hovd</option>
                        <option value="Asia/Jakarta">(+07:00 UTC) Jakarta</option>
                        <option value="Asia/Krasnoyarsk">(+07:00 UTC) Krasnoyarsk</option>
                        <option value="Asia/Novokuznetsk">(+07:00 UTC) Novokuznetsk</option>
                        <option value="Asia/Novosibirsk">(+07:00 UTC) Novosibirsk</option>
                        <option value="Asia/Phnom_Penh">(+07:00 UTC) Phnom Penh</option>
                        <option value="Asia/Pontianak">(+07:00 UTC) Pontianak</option>
                        <option value="Asia/Tomsk">(+07:00 UTC) Tomsk</option>
                        <option value="Asia/Vientiane">(+07:00 UTC) Vientiane</option>
                        <option value="Asia/Brunei">(+08:00 UTC) Brunei</option>
                        <option value="Asia/Choibalsan">(+08:00 UTC) Choibalsan</option>
                        <option value="Asia/Hong_Kong">(+08:00 UTC) Hong Kong</option>
                        <option value="Asia/Irkutsk">(+08:00 UTC) Irkutsk</option>
                        <option value="Asia/Kuala_Lumpur">(+08:00 UTC) Kuala Lumpur</option>
                        <option value="Asia/Kuching">(+08:00 UTC) Kuching</option>
                        <option value="Asia/Macau">(+08:00 UTC) Macau</option>
                        <option value="Asia/Makassar">(+08:00 UTC) Makassar</option>
                        <option value="Asia/Manila">(+08:00 UTC) Manila</option>
                        <option value="Asia/Shanghai">(+08:00 UTC) Shanghai</option>
                        <option value="Asia/Singapore">(+08:00 UTC) Singapore</option>
                        <option value="Asia/Taipei">(+08:00 UTC) Taipei</option>
                        <option value="Asia/Ulaanbaatar">(+08:00 UTC) Ulaanbaatar</option>
                        <option value="Asia/Chita">(+09:00 UTC) Chita</option>
                        <option value="Asia/Dili">(+09:00 UTC) Dili</option>
                        <option value="Asia/Jayapura">(+09:00 UTC) Jayapura</option>
                        <option value="Asia/Khandyga">(+09:00 UTC) Khandyga</option>
                        <option value="Asia/Pyongyang">(+09:00 UTC) Pyongyang</option>
                        <option value="Asia/Seoul">(+09:00 UTC) Seoul</option>
                        <option value="Asia/Tokyo">(+09:00 UTC) Tokyo</option>
                        <option value="Asia/Yakutsk">(+09:00 UTC) Yakutsk</option>
                        <option value="Asia/Ust-Nera">(+10:00 UTC) Ust-Nera</option>
                        <option value="Asia/Vladivostok">(+10:00 UTC) Vladivostok</option>
                        <option value="Asia/Magadan">(+11:00 UTC) Magadan</option>
                        <option value="Asia/Sakhalin">(+11:00 UTC) Sakhalin</option>
                        <option value="Asia/Srednekolymsk">(+11:00 UTC) Srednekolymsk</option>
                        <option value="Asia/Anadyr">(+12:00 UTC) Anadyr</option>
                        <option value="Asia/Kamchatka">(+12:00 UTC) Kamchatka</option>
                    </optgroup>
                    <optgroup label="Atlantic">
                        <option value="Atlantic/Bermuda">(-04:00 UTC) Bermuda</option>
                        <option value="Atlantic/Stanley">(-03:00 UTC) Stanley</option>
                        <option value="Atlantic/South_Georgia">(-02:00 UTC) South Georgia</option>
                        <option value="Atlantic/Azores">(-01:00 UTC) Azores</option>
                        <option value="Atlantic/Cape_Verde">(-01:00 UTC) Cape Verde</option>
                        <option value="Atlantic/Canary">(+00:00 UTC) Canary</option>
                        <option value="Atlantic/Faroe">(+00:00 UTC) Faroe</option>
                        <option value="Atlantic/Madeira">(+00:00 UTC) Madeira</option>
                        <option value="Atlantic/Reykjavik">(+00:00 UTC) Reykjavik</option>
                        <option value="Atlantic/St_Helena">(+00:00 UTC) St Helena</option>
                    </optgroup>
                    <optgroup label="Europe">
                        <option value="Europe/Dublin">(+00:00 UTC) Dublin</option>
                        <option value="Europe/Guernsey">(+00:00 UTC) Guernsey</option>
                        <option value="Europe/Isle_of_Man">(+00:00 UTC) Isle of Man</option>
                        <option value="Europe/Jersey">(+00:00 UTC) Jersey</option>
                        <option value="Europe/Lisbon">(+00:00 UTC) Lisbon</option>
                        <option value="Europe/London">(+00:00 UTC) London</option>
                        <option value="Europe/Amsterdam">(+01:00 UTC) Amsterdam</option>
                        <option value="Europe/Andorra">(+01:00 UTC) Andorra</option>
                        <option value="Europe/Belgrade">(+01:00 UTC) Belgrade</option>
                        <option value="Europe/Berlin">(+01:00 UTC) Berlin</option>
                        <option value="Europe/Bratislava">(+01:00 UTC) Bratislava</option>
                        <option value="Europe/Brussels">(+01:00 UTC) Brussels</option>
                        <option value="Europe/Budapest">(+01:00 UTC) Budapest</option>
                        <option value="Europe/Busingen">(+01:00 UTC) Busingen</option>
                        <option value="Europe/Copenhagen">(+01:00 UTC) Copenhagen</option>
                        <option value="Europe/Gibraltar">(+01:00 UTC) Gibraltar</option>
                        <option value="Europe/Ljubljana">(+01:00 UTC) Ljubljana</option>
                        <option value="Europe/Luxembourg">(+01:00 UTC) Luxembourg</option>
                        <option value="Europe/Madrid">(+01:00 UTC) Madrid</option>
                        <option value="Europe/Malta">(+01:00 UTC) Malta</option>
                        <option value="Europe/Monaco">(+01:00 UTC) Monaco</option>
                        <option value="Europe/Oslo">(+01:00 UTC) Oslo</option>
                        <option value="Europe/Paris">(+01:00 UTC) Paris</option>
                        <option value="Europe/Podgorica">(+01:00 UTC) Podgorica</option>
                        <option value="Europe/Prague">(+01:00 UTC) Prague</option>
                        <option value="Europe/Rome">(+01:00 UTC) Rome</option>
                        <option value="Europe/San_Marino">(+01:00 UTC) San Marino</option>
                        <option value="Europe/Sarajevo">(+01:00 UTC) Sarajevo</option>
                        <option value="Europe/Skopje">(+01:00 UTC) Skopje</option>
                        <option value="Europe/Stockholm">(+01:00 UTC) Stockholm</option>
                        <option value="Europe/Tirane">(+01:00 UTC) Tirane</option>
                        <option value="Europe/Vaduz">(+01:00 UTC) Vaduz</option>
                        <option value="Europe/Vatican">(+01:00 UTC) Vatican</option>
                        <option value="Europe/Vienna">(+01:00 UTC) Vienna</option>
                        <option value="Europe/Warsaw">(+01:00 UTC) Warsaw</option>
                        <option value="Europe/Zagreb">(+01:00 UTC) Zagreb</option>
                        <option value="Europe/Zurich">(+01:00 UTC) Zurich</option>
                        <option value="Europe/Athens">(+02:00 UTC) Athens</option>
                        <option value="Europe/Bucharest">(+02:00 UTC) Bucharest</option>
                        <option value="Europe/Chisinau">(+02:00 UTC) Chisinau</option>
                        <option value="Europe/Helsinki">(+02:00 UTC) Helsinki</option>
                        <option value="Europe/Kaliningrad">(+02:00 UTC) Kaliningrad</option>
                        <option value="Europe/Kiev">(+02:00 UTC) Kiev</option>
                        <option value="Europe/Mariehamn">(+02:00 UTC) Mariehamn</option>
                        <option value="Europe/Riga">(+02:00 UTC) Riga</option>
                        <option value="Europe/Sofia">(+02:00 UTC) Sofia</option>
                        <option value="Europe/Tallinn">(+02:00 UTC) Tallinn</option>
                        <option value="Europe/Uzhgorod">(+02:00 UTC) Uzhgorod</option>
                        <option value="Europe/Vilnius">(+02:00 UTC) Vilnius</option>
                        <option value="Europe/Zaporozhye">(+02:00 UTC) Zaporozhye</option>
                        <option value="Europe/Istanbul">(+03:00 UTC) Istanbul</option>
                        <option value="Europe/Kirov">(+03:00 UTC) Kirov</option>
                        <option value="Europe/Minsk">(+03:00 UTC) Minsk</option>
                        <option value="Europe/Moscow">(+03:00 UTC) Moscow</option>
                        <option value="Europe/Simferopol">(+03:00 UTC) Simferopol</option>
                        <option value="Europe/Astrakhan">(+04:00 UTC) Astrakhan</option>
                        <option value="Europe/Samara">(+04:00 UTC) Samara</option>
                        <option value="Europe/Saratov">(+04:00 UTC) Saratov</option>
                        <option value="Europe/Ulyanovsk">(+04:00 UTC) Ulyanovsk</option>
                        <option value="Europe/Volgograd">(+04:00 UTC) Volgograd</option>
                    </optgroup>
                    <optgroup label="Indian">
                        <option value="Indian/Antananarivo">(+03:00 UTC) Antananarivo</option>
                        <option value="Indian/Comoro">(+03:00 UTC) Comoro</option>
                        <option value="Indian/Mayotte">(+03:00 UTC) Mayotte</option>
                        <option value="Indian/Mahe">(+04:00 UTC) Mahe</option>
                        <option value="Indian/Mauritius">(+04:00 UTC) Mauritius</option>
                        <option value="Indian/Reunion">(+04:00 UTC) Reunion</option>
                        <option value="Indian/Kerguelen">(+05:00 UTC) Kerguelen</option>
                        <option value="Indian/Maldives">(+05:00 UTC) Maldives</option>
                        <option value="Indian/Chagos">(+06:00 UTC) Chagos</option>
                        <option value="Indian/Cocos">(+06:30 UTC) Cocos</option>
                        <option value="Indian/Christmas">(+07:00 UTC) Christmas</option>
                    </optgroup>
                    <optgroup label="Pacific">
                        <option value="Pacific/Midway">(-11:00 UTC) Midway</option>
                        <option value="Pacific/Niue">(-11:00 UTC) Niue</option>
                        <option value="Pacific/Pago_Pago">(-11:00 UTC) Pago Pago</option>
                        <option value="Pacific/Honolulu">(-10:00 UTC) Honolulu</option>
                        <option value="Pacific/Rarotonga">(-10:00 UTC) Rarotonga</option>
                        <option value="Pacific/Tahiti">(-10:00 UTC) Tahiti</option>
                        <option value="Pacific/Gambier">(-09:00 UTC) Gambier</option>
                        <option value="Pacific/Marquesas">(-09:30 UTC) Marquesas</option>
                        <option value="Pacific/Pitcairn">(-08:00 UTC) Pitcairn</option>
                        <option value="Pacific/Galapagos">(-06:00 UTC) Galapagos</option>
                        <option value="Pacific/Easter">(-05:00 UTC) Easter</option>
                        <option value="Pacific/Palau">(+09:00 UTC) Palau</option>
                        <option value="Pacific/Chuuk">(+10:00 UTC) Chuuk</option>
                        <option value="Pacific/Guam">(+10:00 UTC) Guam</option>
                        <option value="Pacific/Port_Moresby">(+10:00 UTC) Port Moresby</option>
                        <option value="Pacific/Saipan">(+10:00 UTC) Saipan</option>
                        <option value="Pacific/Bougainville">(+11:00 UTC) Bougainville</option>
                        <option value="Pacific/Efate">(+11:00 UTC) Efate</option>
                        <option value="Pacific/Guadalcanal">(+11:00 UTC) Guadalcanal</option>
                        <option value="Pacific/Kosrae">(+11:00 UTC) Kosrae</option>
                        <option value="Pacific/Norfolk">(+11:00 UTC) Norfolk</option>
                        <option value="Pacific/Noumea">(+11:00 UTC) Noumea</option>
                        <option value="Pacific/Pohnpei">(+11:00 UTC) Pohnpei</option>
                        <option value="Pacific/Funafuti">(+12:00 UTC) Funafuti</option>
                        <option value="Pacific/Kwajalein">(+12:00 UTC) Kwajalein</option>
                        <option value="Pacific/Majuro">(+12:00 UTC) Majuro</option>
                        <option value="Pacific/Nauru">(+12:00 UTC) Nauru</option>
                        <option value="Pacific/Tarawa">(+12:00 UTC) Tarawa</option>
                        <option value="Pacific/Wake">(+12:00 UTC) Wake</option>
                        <option value="Pacific/Wallis">(+12:00 UTC) Wallis</option>
                        <option value="Pacific/Auckland">(+13:00 UTC) Auckland</option>
                        <option value="Pacific/Enderbury">(+13:00 UTC) Enderbury</option>
                        <option value="Pacific/Fakaofo">(+13:00 UTC) Fakaofo</option>
                        <option value="Pacific/Fiji">(+13:00 UTC) Fiji</option>
                        <option value="Pacific/Tongatapu">(+13:00 UTC) Tongatapu</option>
                        <option value="Pacific/Chatham">(+13:45 UTC) Chatham</option>
                    </optgroup>
                    </select></p>
                <input type="submit" value="Register">
                <a href="/login" style="color: green" class="areadyRegistred">Already registered?</a>
                <h2 style="color:red">Incorrect timezone</h2>
            </form>
        </div>
    </body>`, bodyString)
}

func TestHashAndSaltSuccess(t *testing.T) {
	plainPass := "password"
	hashedPass, err := hashAndSalt(plainPass)
	require.NoError(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass))
	require.NoError(t, err)
}
