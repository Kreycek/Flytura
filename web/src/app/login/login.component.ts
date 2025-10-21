import { Component, ViewChild } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { LoginService } from './login.service';
import { FormsModule } from '@angular/forms';
import { ModalOkComponent } from '../modal/modal-ok/modal-ok.component';
import { ModelsComponent } from '../modulos/companys/models/models/models.component';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [RouterModule, ModalOkComponent, FormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {

  @ViewChild(ModalOkComponent) modalOk!: ModalOkComponent;  
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
             this.router.navigate(['/aplicacao/center']);
          },
          error: (validationError) => {
            // console.error('Erro ao validar o token:', validationError);
            this.errorMessage = 'Erro ao validar o token.';
          }
        });
      },
      error: async (error) => {
        // Verifica se o erro é 401 (Unauthorized)
        if (error.status === 401) {
          console.error('Erro 401: Credenciais inválidas');

           const resultado = await this.modalOk.openModal("Usuário ou senha inválidos",true);             
                      if (resultado) {
                    
                        
                    
                        // Insira aqui a lógica para continuar após a confirmação
                      } else {
                        
                      }
                      
          this.errorMessage = 'Credenciais inválidas. Tente novamente.';
          // Redireciona para a página de login
          // this.router.navigate(['/login']);
        } else {

            const resultado = await this.modalOk.openModal("Sistema com interrupções tente novamente mais tarde.",true);             
                      if (resultado) {
                    
                        
                    
                        // Insira aqui a lógica para continuar após a confirmação
                      } else {
                        
                      }
          // Tratar outros tipos de erro (exemplo: erro de rede, servidor)
          this.errorMessage = error.error?.message || 'Erro desconhecido ao fazer login!';
        }
      }
    });
  }
  
}
  


