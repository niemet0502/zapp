import { CommonModule } from '@angular/common';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { AfterViewInit, Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import * as L from 'leaflet';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule, FormsModule, HttpClientModule],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css',
})
export class DashboardComponent implements AfterViewInit {
  private map!: L.Map;
  marker!: L.Marker;
  searchQuery = '';
  lat!: number;
  lng!: number;

  private eventSource!: EventSource; // SSE connection

  constructor(private http: HttpClient) {}

  ngAfterViewInit(): void {
    this.initMap();
  }

  private initMap(): void {
    this.map = L.map('map', {
      center: [51.505, -0.09], // default center
      zoom: 13,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(this.map);
  }

  searchAddress(): void {
    if (!this.searchQuery) return;

    const url = `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(
      this.searchQuery
    )}`;
    this.http.get<any[]>(url).subscribe((results: any) => {
      if (results.length === 0) {
        alert('Address not found');
        return;
      }

      const { lat, lon } = results[0];
      this.lat = parseFloat(lat);
      this.lng = parseFloat(lon);

      if (this.marker) {
        this.map.removeLayer(this.marker);
      }

      this.marker = L.marker([this.lat, this.lng])
        .addTo(this.map)
        .bindPopup(this.searchQuery)
        .openPopup();

      this.map.setView([this.lat, this.lng], 15);
    });
  }

  requestRide() {
    alert(
      `Ride requested to: ${this.searchQuery}\nLat: ${this.lat}, Lng: ${this.lng}`
    );
    this.startSSE();
  }

  private startSSE(): void {
    const authData = JSON.parse(localStorage.getItem('authData') || '{}');
    console.log(authData.token, authData.email, authData.whatever);

    this.eventSource = new EventSource(
      `http://localhost:3004/sse?type=${authData.userType}&id=${authData.id}`
    );

    this.eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log('Ride update:', data);

      // Example: update driver location on map
      if (data.driverLat && data.driverLng) {
        if (!this.marker) {
          this.marker = L.marker([data.driverLat, data.driverLng]).addTo(
            this.map
          );
        } else {
          this.marker.setLatLng([data.driverLat, data.driverLng]);
        }
      }

      // Example: ride status message
      if (data.status) {
        console.log('Ride status:', data.status);
      }
    };

    this.eventSource.onerror = (err) => {
      console.error('SSE error:', err);
      this.eventSource.close();
    };
  }
}
