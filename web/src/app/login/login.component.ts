import { Component } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { LoginService } from './login.service';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [RouterModule, HttpClientModule, FormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {


  constructor(private loginService: LoginService,private router: Router) {}

  username: string = '';
  email: string = '';
  password: string = '';
  errorMessage: string = '';
  
  onLogin() {
    const user = {
      username: this.username,
      email: this.email,
      password: this.password,
    };
  
    this.loginService.login(user).subscribe({
      next: (response) => {
        // Salva o token no localStorage após o login
        const token = response.token;
        localStorage.setItem('token', token);
  
        // Valida o token após o login
        this.loginService.validateTokenAfterLogin(token).subscribe({
          next: (validationResponse) => {
            this.router.navigate(['/aplicacao']);
          },
          error: (validationError) => {
            // console.error('Erro ao validar o token:', validationError);
            this.errorMessage = 'Erro ao validar o token.';
          }
        });
      },
      error: (error) => {
        // Verifica se o erro é 401 (Unauthorized)
        if (error.status === 401) {
          console.error('Erro 401: Credenciais inválidas');
          this.errorMessage = 'Credenciais inválidas. Tente novamente.';
          // Redireciona para a página de login
          // this.router.navigate(['/login']);
        } else {
          // Tratar outros tipos de erro (exemplo: erro de rede, servidor)
          this.errorMessage = error.error?.message || 'Erro desconhecido ao fazer login!';
        }
      }
    });
  }
  
}
  


