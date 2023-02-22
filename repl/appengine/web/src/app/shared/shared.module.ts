import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TrimPipe } from './trim.pipe';



@NgModule({
  declarations: [
    TrimPipe
  ],
  exports: [
    TrimPipe
  ],
  imports: [
    CommonModule
  ]
})
export class SharedModule { }
