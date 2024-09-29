import axios from "axios";

const instance = axios.create({
	baseURL:   __API_URL__,
	timeout: 1000 * 10
});
instance.interceptors.request.use((config) => {
	console.log("Request URL:", config.url); // Verifica l'URL della richiesta
	return config;
})
export default instance;
