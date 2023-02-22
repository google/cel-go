import { Pipe, PipeTransform } from '@angular/core';

/**
 * Pipe for templates that caps the length of input string (if trimmed, marking
 * with an elipsis).
 */
@Pipe({
  name: 'trim'
})
export class TrimPipe implements PipeTransform {

  transform(value: string, length : number): string {
    return value.length < length ? 
      value :
      value.substring(0, length - 3) + "..."; 
  }

}
