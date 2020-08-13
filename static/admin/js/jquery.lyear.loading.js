;(function(){
    $.fn.lyearloading = function(options) {

        // DOM瀹瑰櫒瀵硅薄
        var $this = $(this);

        var defaults = {
            opacity           : 0.1,                // 閬僵灞傜殑閫忔槑搴︼紝==0鏃舵病鏈夐伄缃╁眰
            backgroundColor   : '#000000',          // 閬僵灞傜殑棰滆壊
            imgUrl            : '',                 // 鍔犺浇鍔ㄧ敾浣跨敤鍥剧墖
            textColorClass    : '',                 // 瀹氫箟鏂囧瓧鐨勯鑹诧紝褰撲笉浣跨敤鍥剧墖鏃�
            spinnerColorClass : '',                 // 瀹氫箟鍔犺浇鍔ㄧ敾鐨勯鑹诧紝褰撲笉浣跨敤鍥剧墖鏃�
            spinnerSize       : 'normal',           // 瀹氫箟鍔犺浇鍔ㄧ敾鐨勫ぇ灏忥紝褰撲笉浣跨敤鍥剧墖鏃�
            spinnerText       : '',                 // 鏄剧ず鐨勬枃瀛�
            zindex            : 9999,               // 閬僵灞傜殑z-index鍊�
        };

        // 铻嶅悎閰嶇疆椤�
        var opts = $.extend({}, defaults, options);

        // 榛樿鏍峰紡
        var maskStyle  = {
            'position'         : 'absolute',
            'width'            : '100%',
            'height'           : '100%',
            'top'              : 0,
            'left'             : 0,
            'background-color' : opts.backgroundColor,
            'opacity'          : opts.opacity,
            'z-index'          : opts.zindex,

        }, textStyle  = {
            'position'       : 'absolute',
            'line-height'    : '120%',
            'text-align'     : 'center',
            'vertical-align' : 'middle',
            'z-index'        : opts.zindex + 1,
        }, spinnerStyle = {
            'position' : 'absolute',
            'z-index'  : opts.zindex + 1,
        };

        var defaultClass = 'lyear-loading';

        // 鍒濆鍖栨柟娉�
        this.init = function(){
            if ($this.children('.' + defaultClass).size() > 0) {
                $this.children('.' + defaultClass).fadeIn(250)
            } else {
                var $maskHtml    = $('<div />').addClass(defaultClass),
                    $textHtml    = $('<span />').html($.trim(opts.spinnerText)).addClass(defaultClass).addClass(opts.textColorClass),
                    $spinnerHtml = opts.imgUrl ? $('<img />').attr('src', opts.imgUrl).addClass(defaultClass) : $('<div />').addClass('spinner-border').addClass(defaultClass).addClass(opts.spinnerColorClass).css(this.getSpinnerSize());

                var toolMethods = {
                    resizeStyle: function() {
                        var $parent        = $this.find('.' + defaultClass).parent(),
                            parentPosition = ('fixed,relative').indexOf($parent.css('position')),
                            isFixed        = parentPosition > -1 || $parent[0] === $('.' + defaultClass)[0].offsetParent,
                            offsetP        = isFixed ? { top: 0, left: 0 } : { top: $parent[0].offsetTop, left: $parent[0].offsetLeft },
                            parentW        = $this.outerWidth(),
                            parentH        = $this.outerHeight();

                        if ($this.selector == 'body') {
                            maskStyle.position     = 'fixed';
                            spinnerStyle.position  = 'fixed';
                            textStyle.position     = 'fixed';

                            spinnerStyle.top  = $(window).height() / 2 - $spinnerHtml.outerHeight() / 2 + (opts.spinnerText ? (- 4 - $textHtml.height() / 2) : 0);
                            spinnerStyle.left = $(window).width() / 2 - $spinnerHtml.outerWidth() / 2;

                            textStyle.top  = $(window).height() / 2 + $spinnerHtml.outerHeight() / 2 - 4;
                            textStyle.left = $(window).width() / 2 - $textHtml.outerWidth() / 2;
                        } else {
                            maskStyle.width  = parentW;
                            maskStyle.height = parentH;
                            maskStyle.top    = offsetP.top;
                            maskStyle.left   = offsetP.left;

                            spinnerStyle.top  = parentH / 2 - $spinnerHtml.outerHeight() / 2 + (opts.spinnerText ? (- 4 - $textHtml.height() / 2) : 0) + offsetP.top;
                            spinnerStyle.left = parentW / 2 - $spinnerHtml.outerWidth() / 2 + offsetP.left;

                            textStyle.top  = parentH / 2 + $spinnerHtml.outerHeight() / 2 - 4 + offsetP.top;
                            textStyle.left = parentW / 2 - $textHtml.width() / 2 + offsetP.left;
                        }

                        $maskHtml.css(maskStyle);
                        $spinnerHtml.css(spinnerStyle);
                        $textHtml.css(textStyle);
                    }
                };

                // 閬僵灞傜户鎵跨埗鍏冪礌鐨勮竟妗嗘晥鏋�
                maskStyle['border-top-left-radius']     = $this.css('border-top-left-radius');
                maskStyle['border-top-right-radius']    = $this.css('border-top-right-radius');
                maskStyle['border-bottom-left-radius']  = $this.css('border-bottom-left-radius');
                maskStyle['border-bottom-right-radius'] = $this.css('border-bottom-right-radius');

                opts.opacity && $maskHtml.css(maskStyle).appendTo($this);
                $.trim(opts.spinnerText) && $textHtml.css(textStyle).appendTo($this);
                $spinnerHtml.css(spinnerStyle).appendTo($this);

                this.loadImage(opts.imgUrl, function (imgObj) {
                    toolMethods.resizeStyle();
                }, function(e){
                    throw new Error(e);
                });

                $(window).off('resize.' + defaultClass).on('resize.' + defaultClass, function () {
                    toolMethods.resizeStyle();
                });
            }
        }

        this.hide = function() {
            $this.children('.' + defaultClass).fadeOut(250);
        }

        this.show = function() {
            $this.children('.' + defaultClass).fadeIn(250);
        }

        this.destroy = function() {
            $this.children('.' + defaultClass).fadeOut(250, function() {
                $(window).off('resize.' + defaultClass);
                $(this).remove();
            });
        }

        this.loadImage = function (url, callback, error) {
            if (!url) {
                return callback();
            }

            var imgObj;

            imgObj     = new Image();
            imgObj.src = url;

            if (imgObj.complete && callback) {
                return callback();
            }

            imgObj.onload = function () {
                imgObj.onload = null;
                callback && callback();
            };

            imgObj.onerror = function (e) {
                imgObj.onerror = null;
                error && error(e);
            };

            return imgObj;
        }

        // 瀵筶oading璁剧疆澶у皬鐨勫鐞嗚繑鍥�
        this.getSpinnerSize = function() {
            var sizeCss;
            switch (options.spinnerSize) {
                case 'sm' :
                    sizeCss = {'width': '12px', 'height' : '12px'};
                    break;
                case 'nm' :
                    sizeCss = {'width': '24px', 'height' : '24px'};
                    break;
                case 'md' :
                    sizeCss = {'width': '36px', 'height' : '36px'};
                    break;
                case 'lg' :
                    sizeCss = {'width': '48px', 'height' : '48px'};
                    break;
                default :
                    sizeCss = {'width': options.spinnerSize, 'height': options.spinnerSize};
            }

            return sizeCss;
        };

        // 鑷姩鎵ц鍒濆鍖栧嚱鏁�
        this.init();

        // 杩斿洖鍑芥暟瀵硅薄
        return this;
    }
})(jQuery);