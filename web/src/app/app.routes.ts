import { RouterModule, Routes } from '@angular/router';
import { CenterComponent } from './center/center.component';
import { NgModule } from '@angular/core';

import { AplicacaoComponent } from './aplicacao/aplicacao.component';
import { LoginComponent } from './login/login.component';
import { UsuarioComponent } from './modulos/ferramentaGestao/users/usuario/usuario.component';
import { AddUsuarioComponent } from './modulos/ferramentaGestao/users/add-usuario/add-usuario.component';


import { InvoicesOnlyFlyComponent } from './modulos/companys/invoices/invoices.component';
import { AddInvoicesComponent } from './modulos/companys/add-invoices/add-invoices.component';

export const routes: Routes = [

    { path: '', redirectTo: '/aplicacao', pathMatch: 'full' },   // Redireciona para login por padrão    
    { 
        path:'aplicacao',
        component:AplicacaoComponent,
        children:[
          
            {
                path:'center',
                component:CenterComponent,
            },
            {
                path:'usuario',
                component:UsuarioComponent,
            }, 
            {
                path:'addUser/:id',
                component:AddUsuarioComponent
            }   , 
            {
                path:'addUser',
                component:AddUsuarioComponent
            },             
            {
                path:'invoicesOnlyFlyExcel',
                component:InvoicesOnlyFlyComponent
            }, 
            {
                path:'addInvoicesOnlyFlyExcel',
                component:AddInvoicesComponent
            }, 
            {
                path:'addInvoicesOnlyFlyExcel/:id',
                component:AddInvoicesComponent
            }

            
        ]
    },   // Redireciona para login por padrão    
    { 
        path:'login',
        component:LoginComponent    
    }
   
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
  })
  export class AppRoutingModule { }