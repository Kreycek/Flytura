import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Component, ViewChild } from '@angular/core';
import { Observable } from 'rxjs';
import { UsuarioService } from '../usuario.service';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { ModalConfirmationComponent } from '../../../../modal/modal-confirmation/modal-confirmation.component';
import { PerfilService } from '../../perfil/perfil.service';

import { FormsModule } from '@angular/forms';
import { PaginatorComponent } from '../../../../paginator/paginator.component';
import { ModalOkComponent } from '../../../../modal/modal-ok/modal-ok.component';
import { ConfigService } from '../../../../services/config.service';


@Component({
  selector: 'app-usuario',
  standalone: true,
  imports: [CommonModule,ModalConfirmationComponent,FormsModule,PaginatorComponent ],
  templateUrl: './usuario.component.html',
  styleUrl: './usuario.component.css'
})
export class UsuarioComponent {

  constructor(
    private router: Router, 
    private usuarioService: UsuarioService,
    private perfilService: PerfilService,
    private configService: ConfigService
  ) {}

   @ViewChild(ModalConfirmationComponent) modal!: ModalConfirmationComponent;
  searchName: string = '';
  searchEmail: string = '';
  searchRole: string = '';

  totalUsers: number = 0;
  totalPages: number = 1;
  currentPage: number = 1;
  limit: number = 0
  filteredUsers = []; // Inicialmente, exibe todos os usuários
  dados:any
  perfis:any
  
  
  ngOnInit() {
    this.perfilService.gePerfil().subscribe((response:any)=>{
      this.perfis=response;
      
    })

    this.limit=this.configService.limitPaginator;


    this.usuarioService.getUsers(this.currentPage,this.limit).subscribe((response:any)=>{
      this.dados=response.users;
      this.totalUsers = response.total;
      this.totalPages = response.pages;
    })
   }
  


  searchUsers(currentPagepage:number) {

    let objPesquisar: { 
      name: string;
      email: string;      
      perfil: number[]; // Definindo o tipo correto para o array 'perfil'
      page:number;
      limit:number;
    }

    objPesquisar= { 
      name: this.searchName, 
      email: this.searchEmail, 
      perfil: [],
      page:currentPagepage,
      limit:this.limit
    };

    if(this.searchRole=='Administrador')
      objPesquisar.perfil.push(1)
    else  if(this.searchRole=='Super Administrador') {
      objPesquisar.perfil.push(2)
    }
    else  if(this.searchRole=='Utilizador') {
      objPesquisar.perfil.push(3)
    }

    this.usuarioService.searchUsers(objPesquisar).subscribe((response:any)=>{

      this.dados=response.users;    
      this.totalUsers = response.total;
      this.totalPages = response.pages;
    })

  }

  async openModal() {
    const resultado = await this.modal.openModal(true,""); 
    if (resultado) {
      // this.modal.isVisible=false;
      this.modal.isVisible=false;
    } else {
      this.modal.isVisible=false;
   
    }
  }

   addUser() {
    this.router.navigate(['/aplicacao/addUser']);
   }

   updateUser(id:string) {
    this.router.navigate(['/aplicacao/addUser', id]);
   
   } 
  
   onPageChange(newPage: number) {
    this.currentPage = newPage;
    this.searchUsers(this.currentPage);
  }
}
