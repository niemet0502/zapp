import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  imports: [CommonModule, FormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css',
})
export class LoginComponent {
  email = '';
  password = '';
  errorMessage = '';
  loading = false;

  constructor(private auth: AuthService) {}

  onLogin() {
    this.loading = true;
    this.errorMessage = '';

    this.auth.login(this.email, this.password).subscribe({
      next: () => (this.loading = false),
      error: (err) => {
        this.loading = false;
        this.errorMessage = err.error?.message || 'Login failed';
      },
    });
  }
}
