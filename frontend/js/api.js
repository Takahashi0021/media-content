const API_BASE = 'http://localhost:8080/api';

class ApiClient {
    constructor() {
        this.token = localStorage.getItem('token');
    }

    setToken(token) {
        this.token = token;
        if (token) {
            localStorage.setItem('token', token);
        } else {
            localStorage.removeItem('token');
        }
    }

    getHeaders() {
        const headers = {
            'Content-Type': 'application/json'
        };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }
        return headers;
    }

    async request(method, endpoint, body = null) {
        const url = `${API_BASE}${endpoint}`;
        
        const options = {
            method,
            headers: this.getHeaders()
        };
        
        if (body) {
            options.body = JSON.stringify(body);
        }

        try {
            const response = await fetch(url, options);
            
            // Парсим JSON ответ
            let data;
            try {
                data = await response.json();
            } catch (e) {
                data = {};
            }
            
            if (!response.ok) {
                throw new Error(data.error || `HTTP ${response.status}`);
            }
            
            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    async register(username, email, password) {
        const result = await this.request('POST', '/auth/register', { username, email, password });
        // Если при регистрации приходит токен - сохраняем
        if (result.token) {
            this.setToken(result.token);
        }
        return result;
    }

    async login(email, password) {
        const result = await this.request('POST', '/auth/login', { email, password });
        if (result.token) {
            this.setToken(result.token);
        }
        return result;
    }

    async logout() {
        this.setToken(null);
        window.location.href = '/login.html';
    }

    async getCurrentUser() {
        return this.request('GET', '/user');
    }

    async getMovies() {
        const result = await this.request('GET', '/movies');
        return result.data || result;
    }

    async getMovie(id) {
        const result = await this.request('GET', `/movies/${id}`);
        return result.data || result;
    }

    async createMovie(movie) {
        return this.request('POST', '/movies', movie);
    }

    async updateMovie(id, data) {
        return this.request('PUT', `/movies/${id}`, data);
    }

    async deleteMovie(id) {
        return this.request('DELETE', `/movies/${id}`);
    }

    async getSeries() {
        const result = await this.request('GET', '/series');
        return result.data || result;
    }

    async getSeriesById(id) {
        const result = await this.request('GET', `/series/${id}`);
        return result.data || result;
    }

    async createSeries(series) {
        return this.request('POST', '/series', series);
    }

    async updateSeries(id, data) {
        return this.request('PUT', `/series/${id}`, data);
    }

    async deleteSeries(id) {
        return this.request('DELETE', `/series/${id}`);
    }

    async getFavorites() {
        const result = await this.request('GET', '/favorites');
        return result.data || result;
    }

    async addFavorite(movieId = null, seriesId = null) {
        const body = {};
        if (movieId) body.movie_id = movieId;
        if (seriesId) body.series_id = seriesId;
        return this.request('POST', '/favorites', body);
    }

    async removeFavorite(id) {
        return this.request('DELETE', `/favorites/${id}`);
    }

    async getReviews(movieId = null, seriesId = null, page = 1, limit = 10) {
        let url = `/reviews?page=${page}&limit=${limit}`;
        if (movieId) url += `&movie_id=${movieId}`;
        if (seriesId) url += `&series_id=${seriesId}`;
        const result = await this.request('GET', url);
        return result;
    }

    async createReview(review) {
        return this.request('POST', '/reviews', review);
    }

    async updateReview(id, rating, comment) {
        return this.request('PUT', `/reviews/${id}`, { rating, comment });
    }

    async deleteReview(id) {
        return this.request('DELETE', `/reviews/${id}`);
    }

    async likeReview(id) {
        return this.request('POST', `/reviews/${id}/like`);
    }

    async getAverageRating(movieId = null, seriesId = null) {
        let url = '/reviews/rating/average';
        if (movieId) url += `?movie_id=${movieId}`;
        if (seriesId) url += `?series_id=${seriesId}`;
        return this.request('GET', url);
    }

    async searchExternalMovie(title) {
        return this.request('GET', `/external/movie/${encodeURIComponent(title)}`);
    }

    async searchExternalMovies(query, page = 1) {
        return this.request('GET', `/external/search?q=${encodeURIComponent(query)}&page=${page}`);
    }
}

const api = new ApiClient();