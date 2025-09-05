import { Component, ViewChild } from '@angular/core';
import { PerfilService } from '../../perfil/perfil.service';
import { CommonModule } from '@angular/common';
import { AbstractControl, FormArray, FormBuilder, FormGroup, ReactiveFormsModule, ValidationErrors, ValidatorFn, Validators } from '@angular/forms';
import { UsuarioService } from '../usuario.service';
import { catchError, tap, throwError } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';
import { ModalOkComponent } from "../../../../modal/modal-ok/modal-ok.component";
import { ActivatedRoute, Router } from '@angular/router';


export function passwordMatchValidator(): ValidatorFn {
  return (control: AbstractControl): ValidationErrors | null => {
    const password = control.get('password')?.value;
    const confirmPassword = control.get('passwordConfirm')?.value;

    return password && confirmPassword && password === confirmPassword
      ? null // válido, os dois campos são iguais
      : { passwordMismatch: true }; // inválido, os campos não são iguais
  };
}

@Component({
  selector: 'app-add-usuario',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, ModalOkComponent],
  templateUrl: './add-usuario.component.html',
  styleUrl: './add-usuario.component.css'
})

export class AddUsuarioComponent {

   validaPerfilSelecionado: ValidatorFn = (control: AbstractControl): ValidationErrors | null => {
    if (control instanceof FormArray) {
      const hasSelected = control.controls.some(ctrl => ctrl.get('value')?.value);     
      return hasSelected ? null :  { perfilNaoSelecionado: true };
    }
    return null;
  };
  
  @ViewChild(ModalOkComponent) modal!: ModalOkComponent;
  
  isEdit=false;
  idUser:string |null = null 
  formulario: FormGroup| null = null;

  constructor(
      private perfilService: PerfilService,
      private fb: FormBuilder,
      private usuarioService: UsuarioService,
      private route: ActivatedRoute,
      private router: Router, 
  ) {} 

  ngOnInit() {

    this.route.paramMap.subscribe(params => {
      const id = params.get('id');  // Substitua 'id' pelo nome do parâmetro

      if(id) {
        this.isEdit=true;
        this.usuarioService.getUserById(id??'0').subscribe((response)=>{
          
          this.idUser=id;    
          this. createFormUser(response);
          this.loadPerfil(response.Perfil);
        })
      }
      else {
        this.isEdit=false;
        this. createFormUser({Active:true});
        this.loadPerfil([]);
      }
    });    
  } 

  loadPerfil(perfisIds:[]) {
      this.perfilService.gePerfil().subscribe((perfis:any)=>{     
        const _perfilForm = this.formulario?.get('perfis') as FormArray;

        perfis.forEach((element:any) => {      
          let perfil=0;
          perfil=perfisIds.filter((perfilId:number)=> {           
              return this.retornaIdTipoPerfil(element.name)==perfilId;            
            })[0];
        
          _perfilForm.push(this.criarPerfil(element.ID,element.name, perfil>0 ? true : false));
        });            
      })  
    }
  
  createFormUser(obj:any) {
      this.formulario = this.fb.group({
        active: [obj.Active, Validators.required],
        nome: [obj.Name, Validators.required],
        sobrenome: [obj.LastName, Validators.required],
        email: [{value:obj.Email, disabled: this.isEdit}, [Validators.required, Validators.email]],
        password: [''],
        passwordConfirm: [''],
        passaporte: [obj.PassportNumber],
        mobile: [obj.Mobile],
        perfis: this.fb.array([], [Validators.required,this.validaPerfilSelecionado] )  // Array para perfis selecionados
        // documentos: this.fb.array([this.criarDocumento()])  // Começa com um documento
      }, 
      // { validators: passwordMatchValidator() }
    ); 
      this.togglePasswordValidation()
    }

  get perfisForm() {
    return (this.formulario?.get('perfis') as FormArray);
  }

  criarPerfil(id:string, name:string, value:boolean): FormGroup {
    return this.fb.group({
      id: [id],
      name: [name],
      value: [value]
    });
  }

  retornaIdTipoPerfil(namePerfil:string ) : number  {
    if(namePerfil.toLowerCase()=='administrador')
      return 1
    else  if(namePerfil.toLowerCase()=="super administrador") {
      return 2
    }
    else if(namePerfil.toLowerCase()=='utilizador') {
      return 3
    }
    else {
      return 0
    }
  }

 retornaNameIdPerfil(perfilId:number ) : string | null {
    if(perfilId==1)
      return 'Administrador'
    else   if(perfilId==2) {
      return 'Super administrador'
    }
    else  if(perfilId==3)  {
      return 'Utilizador'
    }
    else {
      return null
    }
 } 

  gravar() {
   
    
    if (this.formulario?.invalid) {
      this.formulario.markAllAsTouched();
      return;
    }
    const formValues=this.formulario?.value;
    const objGravar: { 
      id?:string |null;
      name: string;
      lastName: string;
      email: string;
      passportNumber: string;
      active:boolean;
      perfil: number[]; // Definindo o tipo correto para o array 'perfil'     
      password:string;
      mobile:string
    } ={
      id:null,
      name:formValues.nome,
      lastName:formValues.sobrenome??'',
      email:formValues.email,
      passportNumber:formValues.passaporte??'',
      active:formValues.active,
      perfil:[],      
      password:formValues.password,
      mobile:formValues.mobile
    }
    objGravar.perfil=[];

    if(formValues.perfis && formValues.perfis.length) {
      formValues.perfis.forEach((element:any) => {
      
          if(element.value) {
            objGravar.perfil.push(this.retornaIdTipoPerfil(element.name))           
          }
      });
    }

    if(this.idUser) {

      objGravar.id=this.idUser
      this.usuarioService.updateUser(objGravar).pipe(
        tap(async (response) =>    {
          this.formulario?.controls['password'].clearValidators();
          this.formulario?.controls['passwordConfirm']?.clearValidators();
          this.formulario?.controls['password'].setValue('',{ emitEvent: false })
          this.formulario?.controls['passwordConfirm'].setValue('',{ emitEvent: false })
          const resultado = await this.modal.openModal(response.message,true); 
          if (resultado) {           
          } else {
       
          }     
        }),
        catchError(async (error: HttpErrorResponse) => {
            
              if (error.status === 500) {     
                const resultado = await this.modal.openModal(error.message,true);
                if (resultado) {
                 
                } else {
              
                }

              }

              if (error.status === 401) {
            
                  // router.navigate(['/login']); // Redireciona para a página de login
              }
              return throwError(() => error);
          })
        
        ).subscribe(()=>{})
    }

    else {

    this.usuarioService.verifyExistsUsers({email:objGravar.email}).pipe(

      tap(async (response)=>{
        if(response.message) {          
          const resultado = await this.modal.openModal("Usuário já cadastrado com esse email",true);
          if (resultado) {           
          } else {

          }         
          
        } else {
              this.usuarioService.addUsers(objGravar).pipe(
                  tap(async () =>{
                  
                    const resultado = await this.modal.openModal("Usuário cadastrado com sucesso",true); 

                    if (resultado) {

                    } else {
                    
                    }                   
                  }),
                  catchError((error: HttpErrorResponse) => {
                      
                        if (error.status === 401) {
                                
                          }
                      return throwError(() => error);
                    })).subscribe(()=>{})              
        }})).subscribe()

    }       
  }
  
  togglePasswordValidation() {
    const passwordControl = this.formulario?.get('password');
    const passwordConfirmControl = this.formulario?.get('passwordConfirm');

    if (this.isEdit) {
      // Se for edição, remove as validações de senha
      passwordControl?.clearValidators();
      passwordConfirmControl?.clearValidators();
    } else {
      // Se for criação, adiciona as validações obrigatórias
      passwordControl?.setValidators([Validators.required, Validators.minLength(6)]);
      passwordConfirmControl?.setValidators([Validators.required, this.matchPasswords()]);
    }

    // Se o usuário digitar uma senha, ativa a validação
    passwordControl?.valueChanges.subscribe(value => {
      if (value) {
        passwordControl?.setValidators([Validators.required, Validators.minLength(6)]);
        passwordConfirmControl?.setValidators([Validators.required, this.matchPasswords()]);
      } else {
        passwordControl?.clearValidators();
        passwordConfirmControl?.clearValidators();
      }
      passwordControl?.updateValueAndValidity();
      passwordConfirmControl?.updateValueAndValidity();
    });
  }

  // Função para verificar se os campos de senha coincidem
  matchPasswords() {
    return (control: any) => {
      const password = this.formulario?.get('password')?.value;
      return control.value === password ? null : { passwordsMismatch: true };
    };
  }
 
  cancel() {
    this.router.navigate(['/aplicacao/usuario']);
  }

}
