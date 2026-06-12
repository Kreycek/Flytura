import { ChangeDetectorRef, Component, ViewChild } from '@angular/core';
import { FormArray, FormBuilder, FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';

import { ActivatedRoute, Router } from '@angular/router';
import { ConfigService } from '../../../../services/config.service';
import { ModuloService } from '../../../modulo.service';
import { ModalOkComponent } from '../../../../modal/modal-ok/modal-ok.component';
import { catchError, tap, throwError } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { PurcharseRecordService } from '../purcharse-record.service';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AirLineService } from '../../../gestion/airLine/airLIne.service';
import { jwtDecode } from 'jwt-decode';
import { EnumPerfil } from '../../../Enum/perfil';

@Component({
  selector: 'app-add-purcharse-record',
  imports: [CommonModule, ReactiveFormsModule, FormsModule,TranslateModule,ModalOkComponent, TranslateModule,],
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

      queryStringKey:string |null = null
      queryStringNombre:string |null = null
      queryStringApellido:string |null = null
      queryStringCompanyCodeAirline:string |null = null
      queryStringDtStart:string |null = null
      queryStringDtEnd:string |null = null
      queryStringStatus:string |null = null
      totalCharPermit:number=0
      keyOrigin:string=''
      statusImportData:any[]

      constructor(     
           private fb: FormBuilder,
           private purcharseRecordService:PurcharseRecordService,
           private route: ActivatedRoute,
           private router: Router, 
           public configService:ConfigService,
           public moduloService:ModuloService,
           private cdr: ChangeDetectorRef,
           private airLineService: AirLineService,
           private translate: TranslateService,
       ) {} 
  
  
       decoded:any
       
    ngOnInit() {


      const token = localStorage.getItem('token'); // ou onde você armazenou o JWT  

          if (token) {
            this.decoded = jwtDecode<any>(token);           
          }

          this.route.paramMap.subscribe(params => {
              const id = params.get('id');  // Substitua 'id' pelo nome do parâmetro             
              this.queryStringKey=params?.get('key')
              this.queryStringNombre=params.get('nombre')
              this.queryStringApellido=params.get('apellido')
              this.queryStringCompanyCodeAirline=params.get('companyCodAirline')
              this.queryStringDtStart=params.get('dtStart')
              this.queryStringDtEnd=params.get('dtEnd')
              this.queryStringStatus=params.get('status')

              
              this.moduloService.getAllStatusImportData().subscribe((response:any)=>{     
                this.statusImportData=response.filter((rsi:any)=>{
                  return rsi.name!='Concluído'
                });

              this.statusImportData=this.configService.sortByKey(this.statusImportData,'name');  

              if(id) {           
                
                this.isEdit=true;

                if(this.statusImportData) {

                this.purcharseRecordService.getPurchaseRecordDataById(id??'0').subscribe((response:any)=>{  

                   const si=this.statusImportData.filter((si:any)=>{
                    return si.name==response.Status
                   })[0]

                   

                    this.keyOrigin=response.Key;
                    this.fileName=response.FileName;  
                    this.id=id;  
                    //Se tiver um código e esse estiver preenchido e for difente de Erro==0002 aplica o código do banco de dados senão aplica o código de fila que é 0001  
                    response.Status=(si &&  (si.code && si.code!='0002')) ? si.code : '0001'
                    this.createForm(response);
                    this.airLineService.getAirLineByCode(response.CompanyCode).subscribe((response)=>{             
                      this.totalCharPermit=response.QtdMinCharKey;
                    });
                })

              }
              }
              else {
                this.isEdit=false;
                this.createForm({Active:true, CompanyCode:''});   
              }          
                });

          });    

          this.airLineService.getAllAirLine().subscribe((response)=>{
          this.airLInes=this.configService.sortByKey(response,'name');                 
      });           
    }  
      
    createForm(obj:any) {
          this.formulario = this.fb.group({
            active: [obj.Active, Validators.required],
            companyCode: [obj.CompanyCode, Validators.required],
            key: [obj.Key, Validators.required],
            name: [obj.Name, Validators.required],
            lastName: [obj.LastName, Validators.required],          
            status: [obj.Status??'', Validators.required],          
          });              
      }    
      
      
    updateCall(objGravar:any) {
         objGravar.id=this.id  
            this.purcharseRecordService.updatePurchaseRecordData(objGravar).pipe(
            tap(async (response:any) => {                
           
              const resultado = await this.modal.openModal(this.translate.instant('Ecra.' + response.message),true); 
              if (resultado) {
                this.redirectSearch();
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
   
  
    async gravar() {    

      if (this.formulario?.invalid) {
        this.formulario.markAllAsTouched();
        return;
      }

      
      const formValues=this.formulario?.value;
      const companyData=this.airLInes.filter((response:any)=>{
        return response.code===formValues.companyCode
      })[0];    

      if(this.totalCharPermit>0 && formValues.key.length<this.totalCharPermit && !this.id) {
            const resultado = await this.modal.openModal(
              this.translate.instant('Ecra.keyMinCharLength') +
                this.totalCharPermit.toString() + 
               this.translate.instant('Ecra.char'),true); 
                if (resultado) {
                  return false
                }   
      }


      //Traz todos os status de importação
       const si=this.statusImportData.filter((si:any)=>{
                  return si.code==formValues.status
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
          status:string,
          idUserInserted:string,
          nameUserInserted:string,
          idUserUpdate:string
          
        } ={
          id:null,
          key:formValues.key,
          name:formValues.name??'',       
          lastName:formValues.lastName??'',       
          active:formValues.active,
          companyCode:formValues.companyCode,
          companyName:companyData ? companyData.name : '',
          fileName:this.fileName ? this.fileName : !this.id ? 'Criado Manualmente' : 'Alterado Manualmente',
          status:si.name,
          idUserInserted: '',
          nameUserInserted:'',
          idUserUpdate:''
          
        }     
   
        objGravar.idUserUpdate= this.decoded.idUser
           objGravar.nameUserInserted= this.decoded.name
        if(this.id) {   
          

          if(this.keyOrigin!=formValues.key) {                         
                this.purcharseRecordService.verifyExistPurchaseRecordData({key:objGravar.key}).subscribe((async (response:any)=>{
                if(response.message) {              
                    const resultado = await this.modal.openModal( this.translate.instant('Ecra.existValueFieldKey'),true); 
                    if (resultado) {
                    
                    }
                } else {
                    this.updateCall(objGravar);
                }          
            }));

          } else {
            this.updateCall(objGravar);
          }           
        }  
        else {
         
          this.purcharseRecordService.verifyExistPurchaseRecordData({key:objGravar.key}).subscribe((async (response:any)=>{
            if(response.message) {              
                const resultado = await this.modal.openModal( this.translate.instant('Ecra.existValueFieldKey'),true); 
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
                const resultado = await this.modal.openModal(this.translate.instant('Ecra.purcharseRecordAddSuccess'),true);             
                if (resultado) {
                  this.redirectSearch();
                  // Insira aqui a lógica para continuar após a confirmação
                } else {
                  
                }
              });
              
              }
          }))    
        
          }       

          return true
      }         

      redirectSearch() {
          this.router.navigate(['/aplicacao/purcharseRecord',
          !this.queryStringKey ? '' : this.queryStringKey,
          !this.queryStringNombre? '' : this.queryStringNombre,
          !this.queryStringApellido? '' : this.queryStringApellido,
          !this.queryStringCompanyCodeAirline? '' : this.queryStringCompanyCodeAirline,
          !this.queryStringDtStart? '' : this.queryStringDtStart,
          !this.queryStringDtEnd? '' : this.queryStringDtEnd,
          !this.queryStringStatus? '' : this.queryStringStatus]);
      }
           
      cancel() {
        this.redirectSearch();
      }
  
      deleteDocument(index:number) {
        this.documentForm.removeAt(index);
      }
  
}
