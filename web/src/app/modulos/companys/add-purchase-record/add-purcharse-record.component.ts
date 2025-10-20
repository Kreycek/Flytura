import { ChangeDetectorRef, Component, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';

import { ActivatedRoute, Router } from '@angular/router';
import { ConfigService } from '../../../services/config.service';
import { ModuloService } from '../../modulo.service';
import { ModalOkComponent } from '../../../modal/modal-ok/modal-ok.component';
import { catchError, tap, throwError } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { PurcharseRecordService } from '../purcharseRecord/purcharse-record.service';
import { TranslateModule } from '@ngx-translate/core';
import { AirLineService } from '../airLine/airLIne.service';

@Component({
  selector: 'app-add-purcharse-record',
  imports: [CommonModule, ReactiveFormsModule, FormsModule,TranslateModule,ModalOkComponent],
  templateUrl: './add-purcharse-record.component.html',
  styleUrl: './add-purcharse-record.component.css'
})
export class AddPurchaseRecordComponent {
  
      @ViewChild(ModalOkComponent) modal!: ModalOkComponent;  
  
        documentMiniFormCod:string=''
        documentMiniFormDescription:string=''
         airLInes:any[]=[]
         fileName:string=''
 
       
       isEdit=false;
       id:string |null = null 
       formulario: FormGroup| null = null;
       years:number[]=[]
       currentYear: number = new Date().getFullYear();
      get documentForm() {
         return (this.formulario?.get('documents') as FormArray);
       }
       constructor(     
           private fb: FormBuilder,
           private purcharseRecordService:PurcharseRecordService,
           private route: ActivatedRoute,
           private router: Router, 
           public configService:ConfigService,
           public moduloService:ModuloService,
           private cdr: ChangeDetectorRef,
           private airLineService: AirLineService,
       ) {} 
  
  
    ngOnInit() {
        this.route.paramMap.subscribe(params => {
          const id = params.get('id');  // Substitua 'id' pelo nome do parâmetro

  
          if(id) {
            this.isEdit=true;
            this.purcharseRecordService.getPurchaseRecordDataById(id??'0').subscribe((response:any)=>{       
              console.log('dados',response);

              this.fileName=response.FileName;
  
              this.id=id;    
              this.createForm(response);   
         
            })
          }
          else {
            this.isEdit=false;
            this.createForm({Active:true, CompanyCode:''});   
          }
       
        });    

             this.airLineService.getAllAirLine().subscribe((response)=>{
          this.airLInes=response;
         

        });
    }
  
      
    createForm(obj:any) {
          this.formulario = this.fb.group({
            active: [obj.Active, Validators.required],
            companyCode: [obj.CompanyCode, Validators.required],
            key: [obj.Key, Validators.required],
            name: [obj.Name, Validators.required],
            lastName: [obj.LastName, Validators.required],          
          });              
      }         
   
  
    gravar() {    

         if (this.formulario?.invalid) {
        this.formulario.markAllAsTouched();
        return;
      }

        const formValues=this.formulario?.value;

        const companyData=this.airLInes.filter((response:any)=>{

          return response.code===formValues.companyCode
        })[0]
        const objGravar: { 
          id?:string |null;
          key: string;
          name: string;
          lastName: string;
          active:boolean;
          companyCode:string,
          companyName:string,
          fileName:string,
          status:string
          
        } ={
          id:null,
          key:formValues.key,
          name:formValues.name??'',       
          lastName:formValues.lastName??'',       
          active:formValues.active,
          companyCode:formValues.companyCode,
          companyName:companyData ? companyData.name : '',
          fileName:this.fileName ? this.fileName : 'Criado Manualmente',
          status:companyData.status ? companyData.status : 'Fila'
          
        }     
   
        if(this.id) {    
            objGravar.id=this.id

  
            this.purcharseRecordService.updatePurchaseRecordData(objGravar).pipe(
            tap(async (response:any) =>    {    

            
            
              const resultado = await this.modal.openModal(response.message,true); 
              if (resultado) {
    
              }
    
            }),
            catchError(async (error: HttpErrorResponse) => {
                
                  if (error.status === 500) {            
                  
                    const resultado = await this.modal.openModal(error.message,true); 
                    if (resultado) {
    
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
    
          this.purcharseRecordService.verifyExistPurchaseRecordData({codDaily:objGravar.key}).subscribe((async (response:any)=>{
            if(response.message) {              
                const resultado = await this.modal.openModal("Esse código de diário já está cadastrado tente outro",true); 
                if (resultado) {
                
                }
            }
            else {
              
              this.purcharseRecordService.addPurchaseRecordData(objGravar).pipe(
                catchError((error: HttpErrorResponse) => {   
                  if (error.status === 401) {
                    ;
                  }
                  return throwError(() => error);
                })
              ).subscribe(async () => {            
              
                // Aguarda o resultado do modal antes de continuar
                const resultado = await this.modal.openModal("Diário cadastrado com sucesso",true);             
                if (resultado) {
              
                  // Insira aqui a lógica para continuar após a confirmação
                } else {
                  
                }
              });
              
              }
          }))
    
        
          }       
      }         
           
      cancel() {
        this.router.navigate(['/aplicacao/purcharseRecord']);
      }
  
      deleteDocument(index:number) {
        this.documentForm.removeAt(index);
      }
  
}
